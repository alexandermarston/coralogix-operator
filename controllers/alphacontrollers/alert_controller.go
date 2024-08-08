/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package alphacontrollers

import (
	"context"
	stdErr "errors"
	"fmt"
	"time"

	webhooks "github.com/coralogix/coralogix-operator/controllers/clientset/grpc/outbound-webhooks"
	"github.com/go-logr/logr"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	utils "github.com/coralogix/coralogix-operator/apis"
	coralogixv1alpha1 "github.com/coralogix/coralogix-operator/apis/coralogix/v1alpha1"
	"github.com/coralogix/coralogix-operator/controllers/clientset"
	alerts "github.com/coralogix/coralogix-operator/controllers/clientset/grpc/alerts/v2"
)

var (
	alertProtoSeverityToSchemaSeverity                               = utils.ReverseMap(coralogixv1alpha1.AlertSchemaSeverityToProtoSeverity)
	alertProtoDayToSchemaDay                                         = utils.ReverseMap(coralogixv1alpha1.AlertSchemaDayToProtoDay)
	alertProtoTimeWindowToSchemaTimeWindow                           = utils.ReverseMap(coralogixv1alpha1.AlertSchemaTimeWindowToProtoTimeWindow)
	alertProtoAutoRetireRatioToSchemaAutoRetireRatio                 = utils.ReverseMap(coralogixv1alpha1.AlertSchemaAutoRetireRatioToProtoAutoRetireRatio)
	alertProtoFiltersLogSeverityToSchemaFiltersLogSeverity           = utils.ReverseMap(coralogixv1alpha1.AlertSchemaFiltersLogSeverityToProtoFiltersLogSeverity)
	alertProtoRelativeTimeFrameToSchemaTimeFrameAndRelativeTimeFrame = utils.ReverseMap(coralogixv1alpha1.AlertSchemaRelativeTimeFrameToProtoTimeFrameAndRelativeTimeFrame)
	alertProtoArithmeticOperatorToSchemaArithmeticOperator           = utils.ReverseMap(coralogixv1alpha1.AlertSchemaArithmeticOperatorToProtoArithmeticOperator)
	alertProtoNotifyOn                                               = utils.ReverseMap(coralogixv1alpha1.AlertSchemaNotifyOnToProtoNotifyOn)
	alertProtoFlowOperatorToProtoFlowOperator                        = utils.ReverseMap(coralogixv1alpha1.AlertSchemaFlowOperatorToProtoFlowOperator)
	alertFinalizerName                                               = "alert.coralogix.com/finalizer"
)

// AlertReconciler reconciles a Alert object
type AlertReconciler struct {
	client.Client
	CoralogixClientSet clientset.ClientSetInterface
	Scheme             *runtime.Scheme
}

//+kubebuilder:rbac:groups=coralogix.com,resources=alerts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=coralogix.com,resources=alerts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=coralogix.com,resources=alerts/finalizers,verbs=update

func (r *AlertReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	var (
		err error
	)
	log := log.FromContext(ctx).WithValues(
		"alert", req.NamespacedName.Name,
		"namespace", req.NamespacedName.Namespace,
	)

	log.V(1).Info("Reconciling Alert")
	coralogixv1alpha1.WebhooksClient = r.CoralogixClientSet.OutboundWebhooks()
	alert := coralogixv1alpha1.NewAlert()

	if err = r.Client.Get(ctx, req.NamespacedName, alert); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	if ptr.Deref(alert.Status.ID, "") == "" {
		err = r.create(ctx, log, alert)
		if err != nil {
			log.Error(err, "Error on creating alert")
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
		return ctrl.Result{}, nil
	}

	if !alert.ObjectMeta.DeletionTimestamp.IsZero() {
		err = r.delete(ctx, log, alert)
		if err != nil {
			log.Error(err, "Error on deleting alert")
			return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
		}
		return ctrl.Result{}, nil
	}

	err = r.update(ctx, log, alert)
	if err != nil {
		log.Error(err, "Error on updating alert")
		return ctrl.Result{RequeueAfter: defaultErrRequeuePeriod}, err
	}

	return ctrl.Result{}, nil
}

func (r *AlertReconciler) update(ctx context.Context,
	log logr.Logger,
	alert *coralogixv1alpha1.Alert) error {
	alertRequest, err := alert.Spec.ExtractUpdateAlertRequest(ctx, log, *alert.Status.ID)
	if err != nil {
		return fmt.Errorf("error to parse alert request: %w", err)
	}

	log.V(1).Info("Updating remote alert", "alert", protojson.Format(alertRequest))
	remoteUpdatedAlert, err := r.CoralogixClientSet.Alerts().UpdateAlert(ctx, alertRequest)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			log.Info("alert not found on remote, recreating it")
			alert.Status = *coralogixv1alpha1.NewDefaultAlertStatus()
			if err = r.Status().Update(ctx, alert); err != nil {
				return fmt.Errorf("error on updating alert status: %w", err)
			}
			return fmt.Errorf("alert not found on remote, recreating it: %w", err)
		}
		return fmt.Errorf("error on updating alert: %w", err)
	}
	log.V(1).Info("Remote alert updated", "alert", protojson.Format(remoteUpdatedAlert))

	status, err := getStatus(ctx, log, remoteUpdatedAlert.GetAlert(), alert.Spec)
	if err != nil {
		return fmt.Errorf("error on getting status: %w", err)
	}

	if err = r.Get(ctx, client.ObjectKeyFromObject(alert), alert); err != nil {
		return fmt.Errorf("error on getting alert: %w", err)
	}
	alert.Status = status

	if err = r.Status().Update(ctx, alert); err != nil {
		return fmt.Errorf("error on updating alert status: %w", err)
	}

	return nil
}

func (r *AlertReconciler) delete(ctx context.Context,
	log logr.Logger,
	alert *coralogixv1alpha1.Alert) error {

	log.V(1).Info("Deleting remote alert", "alert", *alert.Status.ID)
	_, err := r.CoralogixClientSet.Alerts().DeleteAlert(ctx, &alerts.DeleteAlertByUniqueIdRequest{
		Id: wrapperspb.String(*alert.Status.ID),
	})
	if err != nil && status.Code(err) != codes.NotFound {
		return fmt.Errorf("error on deleting alert: %w", err)
	}
	log.V(1).Info("Remote alert deleted", "alert", *alert.Status.ID)

	controllerutil.RemoveFinalizer(alert, alertFinalizerName)
	if err = r.Update(ctx, alert); err != nil {
		return fmt.Errorf("error on updating alert: %w", err)
	}

	return nil
}

func (r *AlertReconciler) create(
	ctx context.Context,
	log logr.Logger,
	alert *coralogixv1alpha1.Alert) error {

	if alert.Spec.Labels == nil {
		alert.Spec.Labels = make(map[string]string)
	}

	if value, ok := alert.Spec.Labels["managed-by"]; !ok || value == "" {
		alert.Spec.Labels["managed-by"] = "coralogix-operator"
	}

	if err := r.Update(ctx, alert); err != nil {
		return fmt.Errorf("error on updating alert: %w", err)
	}

	alertRequest, err := alert.ExtractCreateAlertRequest(ctx, log)
	if err != nil {
		return fmt.Errorf("error to parse alert request: %w", err)
	}

	log.V(1).Info("Creating remote alert", "alert", protojson.Format(alertRequest))
	response, err := r.CoralogixClientSet.Alerts().CreateAlert(ctx, alertRequest)
	if err != nil {
		return fmt.Errorf("error on creating alert: %w", err)
	}
	log.V(1).Info("Remote alert created", "response", protojson.Format(response))

	if err = r.Get(ctx, client.ObjectKeyFromObject(alert), alert); err != nil {
		return fmt.Errorf("error on getting alert: %w", err)
	}

	alert.Status.ID = pointer.String(response.GetAlert().GetUniqueIdentifier().GetValue())
	if err = r.Update(ctx, alert); err != nil {
		return fmt.Errorf("error on updating alert: %w", err)
	}

	if alert.Status, err = getStatus(ctx, log, response.GetAlert(), alert.Spec); err != nil {
		return fmt.Errorf("error on getting status: %w", err)
	}
	if err = r.Status().Update(ctx, alert); err != nil {
		return fmt.Errorf("error on updating alert status: %w", err)
	}

	if !controllerutil.ContainsFinalizer(alert, alertFinalizerName) {
		controllerutil.AddFinalizer(alert, alertFinalizerName)
	}
	if err = r.Client.Update(ctx, alert); err != nil {
		return fmt.Errorf("error on updating alert: %w", err)
	}

	return nil
}

func getStatus(ctx context.Context, log logr.Logger, actualAlert *alerts.Alert, spec coralogixv1alpha1.AlertSpec) (coralogixv1alpha1.AlertStatus, error) {
	if actualAlert == nil {
		return coralogixv1alpha1.AlertStatus{}, stdErr.New("alert is nil")
	}

	var status coralogixv1alpha1.AlertStatus
	var err error

	status.ID = utils.WrapperspbStringToStringPointer(actualAlert.GetUniqueIdentifier())

	status.Name = actualAlert.GetName().GetValue()

	status.Description = actualAlert.GetDescription().GetValue()

	status.Active = actualAlert.GetIsActive().GetValue()

	status.Severity = alertProtoSeverityToSchemaSeverity[actualAlert.GetSeverity()]

	status.Labels = flattenMetaLabels(actualAlert.GetMetaLabels())

	status.ExpirationDate = flattenExpirationDate(actualAlert.GetExpiration())

	var timeZone coralogixv1alpha1.TimeZone
	if spec.Scheduling != nil {
		timeZone = spec.Scheduling.TimeZone
	}

	status.Scheduling = flattenScheduling(actualAlert.GetActiveWhen(), timeZone)

	status.AlertType = flattenAlertType(actualAlert)

	if notificationGroups, flattenErr := flattenNotificationGroups(ctx, log, actualAlert.GetNotificationGroups()); flattenErr != nil {
		err = stdErr.Join(err, fmt.Errorf("error on flatten alert - %w", flattenErr))
	} else {
		status.NotificationGroups = notificationGroups
	}

	status.ShowInInsight = flattenShowInInsight(actualAlert.GetShowInInsight())

	status.PayloadFilters = utils.WrappedStringSliceToStringSlice(actualAlert.GetNotificationPayloadFilters())

	return status, err
}

func flattenAlertType(actualAlert *alerts.Alert) coralogixv1alpha1.AlertType {
	actualFilters := actualAlert.GetFilters()
	actualCondition := actualAlert.GetCondition()

	var alertType coralogixv1alpha1.AlertType
	switch actualFilters.GetFilterType() {
	case alerts.AlertFilters_FILTER_TYPE_TEXT_OR_UNSPECIFIED:
		if newValueCondition, ok := actualCondition.GetCondition().(*alerts.AlertCondition_NewValue); ok {
			alertType.NewValue = flattenNewValueAlert(actualFilters, newValueCondition)
		} else {
			alertType.Standard = flattenStandardAlert(actualFilters, actualCondition)
		}
	case alerts.AlertFilters_FILTER_TYPE_RATIO:
		alertType.Ratio = flattenRatioAlert(actualFilters, actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_UNIQUE_COUNT:
		alertType.UniqueCount = flattenUniqueCountAlert(actualFilters, actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_TIME_RELATIVE:
		alertType.TimeRelative = flattenTimeRelativeAlert(actualFilters, actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_METRIC:
		alertType.Metric = flattenMetricAlert(actualFilters, actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_TRACING:
		alertType.Tracing = flattenTracingAlert(actualAlert.GetTracingAlert(), actualCondition)
	case alerts.AlertFilters_FILTER_TYPE_FLOW:
		alertType.Flow = flattenFlowAlert(actualCondition.GetFlow())
	}

	return alertType
}

func flattenNewValueAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition_NewValue) *coralogixv1alpha1.NewValue {
	flattenedFilters := flattenFilters(filters)
	newValueCondition := flattenNewValueCondition(condition.NewValue.GetParameters())

	newValue := &coralogixv1alpha1.NewValue{
		Filters:    flattenedFilters,
		Conditions: newValueCondition,
	}

	return newValue
}

func flattenFilters(filters *alerts.AlertFilters) *coralogixv1alpha1.Filters {
	if filters == nil {
		return nil
	}

	var flattenedFilters = &coralogixv1alpha1.Filters{}

	flattenedFilters.SearchQuery = utils.WrapperspbStringToStringPointer(filters.GetText())

	flattenedFilters.Alias = utils.WrapperspbStringToStringPointer(filters.GetAlias())

	flattenedFilters.Severities = flattenSeverities(filters.GetSeverities())

	if metaData := filters.Metadata; metaData != nil {
		flattenedFilters.Subsystems = utils.WrappedStringSliceToStringSlice(metaData.Subsystems)
		flattenedFilters.Categories = utils.WrappedStringSliceToStringSlice(metaData.Categories)
		flattenedFilters.Applications = utils.WrappedStringSliceToStringSlice(metaData.Applications)
		flattenedFilters.Computers = utils.WrappedStringSliceToStringSlice(metaData.Computers)
		flattenedFilters.Classes = utils.WrappedStringSliceToStringSlice(metaData.Classes)
		flattenedFilters.Methods = utils.WrappedStringSliceToStringSlice(metaData.Methods)
		flattenedFilters.IPs = utils.WrappedStringSliceToStringSlice(metaData.IpAddresses)
	}

	return flattenedFilters
}

func flattenSeverities(severities []alerts.AlertFilters_LogSeverity) []coralogixv1alpha1.FiltersLogSeverity {
	flattenedSeverities := make([]coralogixv1alpha1.FiltersLogSeverity, 0, len(severities))
	for _, severity := range severities {
		sev := alertProtoFiltersLogSeverityToSchemaFiltersLogSeverity[severity]
		flattenedSeverities = append(flattenedSeverities, sev)
	}
	return flattenedSeverities
}

func flattenNewValueCondition(conditionParams *alerts.ConditionParameters) coralogixv1alpha1.NewValueConditions {
	var key string
	if actualKeys := conditionParams.GetGroupBy(); len(actualKeys) != 0 {
		key = actualKeys[0].GetValue()
	}
	timeWindow := coralogixv1alpha1.NewValueTimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])

	newValueCondition := coralogixv1alpha1.NewValueConditions{
		Key:        key,
		TimeWindow: timeWindow,
	}

	return newValueCondition
}

func flattenStandardAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) *coralogixv1alpha1.Standard {
	flattenedFilters := flattenFilters(filters)
	standardCondition := flattenStandardCondition(condition)

	standard := &coralogixv1alpha1.Standard{
		Filters:    flattenedFilters,
		Conditions: standardCondition,
	}

	return standard
}

func flattenStandardCondition(condition *alerts.AlertCondition) coralogixv1alpha1.StandardConditions {
	var standardCondition coralogixv1alpha1.StandardConditions
	var conditionParams *alerts.ConditionParameters

	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		conditionParams = condition.LessThan.GetParameters()
		standardCondition.AlertWhen = coralogixv1alpha1.StandardAlertWhenLessThan
		standardCondition.Threshold = new(int)
		standardCondition.TimeWindow = new(coralogixv1alpha1.TimeWindow)
		*standardCondition.Threshold = int(conditionParams.GetThreshold().GetValue())
		*standardCondition.TimeWindow = coralogixv1alpha1.TimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])

		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			standardCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			standardCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: false,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		standardCondition.AlertWhen = coralogixv1alpha1.StandardAlertWhenMoreThan
		standardCondition.Threshold = new(int)
		standardCondition.TimeWindow = new(coralogixv1alpha1.TimeWindow)
		*standardCondition.Threshold = int(conditionParams.GetThreshold().GetValue())
		*standardCondition.TimeWindow = coralogixv1alpha1.TimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])
	case *alerts.AlertCondition_MoreThanUsual:
		conditionParams = condition.MoreThanUsual.GetParameters()
		standardCondition.AlertWhen = coralogixv1alpha1.StandardAlertWhenMoreThanUsual
		standardCondition.Threshold = new(int)
		*standardCondition.Threshold = int(conditionParams.GetThreshold().GetValue())
	case *alerts.AlertCondition_Immediate:
		standardCondition.AlertWhen = coralogixv1alpha1.StandardAlertWhenImmediately
		return standardCondition
	}

	standardCondition.GroupBy = utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy())

	return standardCondition
}

func flattenRatioAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) *coralogixv1alpha1.Ratio {
	query1Filters := flattenFilters(filters)
	q2Filters := filters.GetRatioAlerts()[0]
	query2Filters := flattenRatioFilters(q2Filters)
	ratioCondition := flattenRatioCondition(condition, q2Filters.GetGroupBy())

	ratio := &coralogixv1alpha1.Ratio{
		Query1Filters: *query1Filters,
		Query2Filters: query2Filters,
		Conditions:    ratioCondition,
	}

	return ratio
}

func flattenRatioFilters(filters *alerts.AlertFilters_RatioAlert) coralogixv1alpha1.RatioQ2Filters {
	var flattenedFilters coralogixv1alpha1.RatioQ2Filters
	if filters == nil {
		return flattenedFilters
	}

	flattenedFilters.SearchQuery = utils.WrapperspbStringToStringPointer(filters.GetText())

	flattenedFilters.Alias = utils.WrapperspbStringToStringPointer(filters.GetAlias())

	flattenedFilters.Severities = flattenSeverities(filters.GetSeverities())
	flattenedFilters.Subsystems = utils.WrappedStringSliceToStringSlice(filters.GetSubsystems())
	flattenedFilters.Applications = utils.WrappedStringSliceToStringSlice(filters.GetApplications())

	return flattenedFilters
}

func flattenRatioCondition(condition *alerts.AlertCondition, groupByQ2 []*wrapperspb.StringValue) coralogixv1alpha1.RatioConditions {
	var ratioCondition coralogixv1alpha1.RatioConditions
	var conditionParams *alerts.ConditionParameters

	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		conditionParams = condition.LessThan.GetParameters()
		ratioCondition.AlertWhen = coralogixv1alpha1.AlertWhenLessThan

		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			ratioCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			ratioCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: false,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		ratioCondition.AlertWhen = coralogixv1alpha1.AlertWhenMoreThan
	}

	ratioCondition.Ratio = utils.FloatToQuantity(conditionParams.GetThreshold().GetValue())
	ratioCondition.TimeWindow = coralogixv1alpha1.TimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])

	if groupByQ1 := conditionParams.GetGroupBy(); len(groupByQ1) > 0 && len(groupByQ2) == 0 {
		ratioCondition.GroupBy = utils.WrappedStringSliceToStringSlice(groupByQ1)
		ratioCondition.GroupByFor = new(coralogixv1alpha1.GroupByFor)
		*ratioCondition.GroupByFor = coralogixv1alpha1.GroupByForQ1
	} else if len(groupByQ2) > 0 && len(groupByQ1) == 0 {
		ratioCondition.GroupBy = utils.WrappedStringSliceToStringSlice(groupByQ2)
		ratioCondition.GroupByFor = new(coralogixv1alpha1.GroupByFor)
		*ratioCondition.GroupByFor = coralogixv1alpha1.GroupByForQ2
	} else if len(groupByQ1) > 0 && len(groupByQ2) > 0 {
		ratioCondition.GroupBy = utils.WrappedStringSliceToStringSlice(groupByQ2)
		ratioCondition.GroupByFor = new(coralogixv1alpha1.GroupByFor)
		*ratioCondition.GroupByFor = coralogixv1alpha1.GroupByForBoth
	}

	return ratioCondition
}

func flattenUniqueCountAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) *coralogixv1alpha1.UniqueCount {
	flattenedFilters := flattenFilters(filters)
	uniqueCountCondition := flattenUniqueCountCondition(condition)

	ratio := &coralogixv1alpha1.UniqueCount{
		Filters:    flattenedFilters,
		Conditions: uniqueCountCondition,
	}

	return ratio
}

func flattenUniqueCountCondition(condition *alerts.AlertCondition) coralogixv1alpha1.UniqueCountConditions {
	conditionParams := condition.GetCondition().(*alerts.AlertCondition_UniqueCount).UniqueCount.GetParameters()
	var uniqueCountCondition coralogixv1alpha1.UniqueCountConditions

	uniqueCountCondition.Key = conditionParams.GetCardinalityFields()[0].GetValue()
	uniqueCountCondition.MaxUniqueValues = int(conditionParams.GetThreshold().GetValue())
	uniqueCountCondition.TimeWindow = coralogixv1alpha1.UniqueValueTimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])
	if actualGroupBy := conditionParams.GetGroupBy(); len(actualGroupBy) > 0 {
		uniqueCountCondition.GroupBy = new(string)
		*uniqueCountCondition.GroupBy = actualGroupBy[0].GetValue()

		uniqueCountCondition.MaxUniqueValuesForGroupBy = new(int)
		*uniqueCountCondition.MaxUniqueValuesForGroupBy = int(conditionParams.GetMaxUniqueCountValuesForGroupByKey().GetValue())
	}

	return uniqueCountCondition
}

func flattenTimeRelativeAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) *coralogixv1alpha1.TimeRelative {
	flattenedFilters := flattenFilters(filters)
	timeRelativeCondition := flattenTimeRelativeCondition(condition)

	timeRelative := &coralogixv1alpha1.TimeRelative{
		Filters:    flattenedFilters,
		Conditions: timeRelativeCondition,
	}

	return timeRelative
}

func flattenTimeRelativeCondition(condition *alerts.AlertCondition) coralogixv1alpha1.TimeRelativeConditions {
	var conditionParams *alerts.ConditionParameters
	var timeRelativeCondition coralogixv1alpha1.TimeRelativeConditions

	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		conditionParams = condition.LessThan.GetParameters()
		timeRelativeCondition.AlertWhen = coralogixv1alpha1.AlertWhenLessThan

		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			timeRelativeCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			timeRelativeCondition.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: false,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		timeRelativeCondition.AlertWhen = coralogixv1alpha1.AlertWhenMoreThan
	}

	timeRelativeCondition.Threshold = utils.FloatToQuantity(conditionParams.GetThreshold().GetValue())
	relativeTimeFrame := coralogixv1alpha1.ProtoTimeFrameAndRelativeTimeFrame{TimeFrame: conditionParams.GetTimeframe(), RelativeTimeFrame: conditionParams.GetRelativeTimeframe()}
	timeRelativeCondition.TimeWindow = alertProtoRelativeTimeFrameToSchemaTimeFrameAndRelativeTimeFrame[relativeTimeFrame]
	timeRelativeCondition.IgnoreInfinity = conditionParams.GetIgnoreInfinity().GetValue()
	timeRelativeCondition.GroupBy = utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy())

	return timeRelativeCondition
}

func flattenMetricAlert(filters *alerts.AlertFilters, condition *alerts.AlertCondition) *coralogixv1alpha1.Metric {
	metric := new(coralogixv1alpha1.Metric)

	var conditionParams *alerts.ConditionParameters
	var promqlAlertWhen coralogixv1alpha1.PromqlAlertWhen
	var luceneAlertWhen coralogixv1alpha1.AlertWhen
	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_LessThan:
		promqlAlertWhen = coralogixv1alpha1.PromqlAlertWhenLessThan
		luceneAlertWhen = coralogixv1alpha1.AlertWhenLessThan
		conditionParams = condition.LessThan.GetParameters()
	case *alerts.AlertCondition_MoreThan:
		conditionParams = condition.MoreThan.GetParameters()
		promqlAlertWhen = coralogixv1alpha1.PromqlAlertWhenMoreThan
		luceneAlertWhen = coralogixv1alpha1.AlertWhenMoreThan
	case *alerts.AlertCondition_MoreThanUsual:
		conditionParams = condition.MoreThanUsual.GetParameters()
		promqlAlertWhen = coralogixv1alpha1.PromqlAlertWhenMoreThanUsual
	}

	if promqlParams := conditionParams.GetMetricAlertPromqlParameters(); promqlParams != nil {
		metric.Promql = flattenPromqlAlert(conditionParams, promqlParams, promqlAlertWhen)
	} else {
		metric.Lucene = flattenLuceneAlert(conditionParams, filters.GetText(), luceneAlertWhen)
	}

	return metric
}

func flattenPromqlAlert(conditionParams *alerts.ConditionParameters, promqlParams *alerts.MetricAlertPromqlConditionParameters, alertWhen coralogixv1alpha1.PromqlAlertWhen) *coralogixv1alpha1.Promql {
	promql := new(coralogixv1alpha1.Promql)

	promql.SearchQuery = promqlParams.GetPromqlText().GetValue()
	promql.Conditions = coralogixv1alpha1.PromqlConditions{
		AlertWhen:                   alertWhen,
		Threshold:                   utils.FloatToQuantity(conditionParams.GetThreshold().GetValue()),
		SampleThresholdPercentage:   int(promqlParams.GetSampleThresholdPercentage().GetValue()),
		TimeWindow:                  coralogixv1alpha1.MetricTimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()]),
		ReplaceMissingValueWithZero: promqlParams.GetSwapNullValues().GetValue(),
	}

	if minNonNullValuesPercentage := promqlParams.GetNonNullPercentage(); minNonNullValuesPercentage != nil {
		promql.Conditions.MinNonNullValuesPercentage = new(int)
		*promql.Conditions.MinNonNullValuesPercentage = int(minNonNullValuesPercentage.GetValue())
	}

	if alertWhen == coralogixv1alpha1.PromqlAlertWhenLessThan {
		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			promql.Conditions.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			promql.Conditions.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: false,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	}

	return promql
}

func flattenLuceneAlert(conditionParams *alerts.ConditionParameters, searchQuery *wrapperspb.StringValue, alertWhen coralogixv1alpha1.AlertWhen) *coralogixv1alpha1.Lucene {
	lucene := new(coralogixv1alpha1.Lucene)
	metricParams := conditionParams.GetMetricAlertParameters()

	if searchQuery != nil {
		lucene.SearchQuery = new(string)
		*lucene.SearchQuery = searchQuery.GetValue()
	}

	lucene.Conditions = coralogixv1alpha1.LuceneConditions{
		MetricField:                 metricParams.GetMetricField().GetValue(),
		ArithmeticOperator:          alertProtoArithmeticOperatorToSchemaArithmeticOperator[metricParams.GetArithmeticOperator()],
		AlertWhen:                   alertWhen,
		Threshold:                   utils.FloatToQuantity(conditionParams.GetThreshold().GetValue()),
		SampleThresholdPercentage:   int(metricParams.GetSampleThresholdPercentage().GetValue()),
		TimeWindow:                  coralogixv1alpha1.MetricTimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()]),
		GroupBy:                     utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy()),
		ReplaceMissingValueWithZero: metricParams.GetSwapNullValues().GetValue(),
		MinNonNullValuesPercentage:  int(metricParams.GetNonNullPercentage().GetValue()),
	}

	if arithmeticOperatorModifier := metricParams.GetArithmeticOperatorModifier(); arithmeticOperatorModifier != nil {
		lucene.Conditions.ArithmeticOperatorModifier = new(int)
		*lucene.Conditions.ArithmeticOperatorModifier = int(arithmeticOperatorModifier.GetValue())
	}

	if alertWhen == coralogixv1alpha1.AlertWhenLessThan {
		if actualManageUndetectedValues := conditionParams.GetRelatedExtendedData(); actualManageUndetectedValues != nil {
			actualShouldTriggerDeadman, actualCleanupDeadmanDuration := actualManageUndetectedValues.GetShouldTriggerDeadman().GetValue(), actualManageUndetectedValues.GetCleanupDeadmanDuration()
			autoRetireRatio := alertProtoAutoRetireRatioToSchemaAutoRetireRatio[actualCleanupDeadmanDuration]
			lucene.Conditions.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: actualShouldTriggerDeadman,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		} else {
			autoRetireRatio := coralogixv1alpha1.AutoRetireRatioNever
			lucene.Conditions.ManageUndetectedValues = &coralogixv1alpha1.ManageUndetectedValues{
				EnableTriggeringOnUndetectedValues: false,
				AutoRetireRatio:                    &autoRetireRatio,
			}
		}
	}

	return lucene
}

func flattenTracingAlert(tracingAlert *alerts.TracingAlert, condition *alerts.AlertCondition) *coralogixv1alpha1.Tracing {
	latencyThresholdMS := float64(tracingAlert.GetConditionLatency()) / float64(time.Millisecond.Microseconds())
	tracingFilters := flattenTracingAlertFilters(tracingAlert)
	tracingFilters.LatencyThresholdMilliseconds = utils.FloatToQuantity(latencyThresholdMS)

	tracingCondition := flattenTracingCondition(condition)

	return &coralogixv1alpha1.Tracing{
		Filters:    tracingFilters,
		Conditions: tracingCondition,
	}
}

func flattenTracingCondition(condition *alerts.AlertCondition) coralogixv1alpha1.TracingCondition {
	var tracingCondition coralogixv1alpha1.TracingCondition
	switch condition := condition.GetCondition().(type) {
	case *alerts.AlertCondition_Immediate:
		tracingCondition.AlertWhen = coralogixv1alpha1.TracingAlertWhenImmediately
	case *alerts.AlertCondition_MoreThan:
		conditionParams := condition.MoreThan.GetParameters()
		tracingCondition.AlertWhen = coralogixv1alpha1.TracingAlertWhenMore

		tracingCondition.Threshold = new(int)
		*tracingCondition.Threshold = int(conditionParams.GetThreshold().GetValue())

		tracingCondition.TimeWindow = new(coralogixv1alpha1.TimeWindow)
		*tracingCondition.TimeWindow = coralogixv1alpha1.TimeWindow(alertProtoTimeWindowToSchemaTimeWindow[conditionParams.GetTimeframe()])

		tracingCondition.GroupBy = utils.WrappedStringSliceToStringSlice(conditionParams.GetGroupBy())
	}

	return tracingCondition
}

func flattenTracingAlertFilters(tracingAlert *alerts.TracingAlert) coralogixv1alpha1.TracingFilters {
	applications, subsystems, services := flattenTracingFilters(tracingAlert.GetFieldFilters())
	tagFilters := flattenTagFiltersData(tracingAlert.GetTagFilters())

	return coralogixv1alpha1.TracingFilters{
		TagFilters:   tagFilters,
		Applications: applications,
		Subsystems:   subsystems,
		Services:     services,
	}
}

func flattenTracingFilters(tracingFilters []*alerts.FilterData) (applications, subsystems, services []string) {
	filtersData := flattenFiltersData(tracingFilters)
	applications = filtersData["applicationName"]
	subsystems = filtersData["subsystemName"]
	services = filtersData["serviceName"]
	return
}

func flattenTagFiltersData(filtersData []*alerts.FilterData) []coralogixv1alpha1.TagFilter {
	fieldToFilters := flattenFiltersData(filtersData)
	result := make([]coralogixv1alpha1.TagFilter, 0, len(fieldToFilters))
	for field, filters := range fieldToFilters {
		filterSchema := coralogixv1alpha1.TagFilter{
			Field:  field,
			Values: filters,
		}
		result = append(result, filterSchema)
	}
	return result
}

func flattenFiltersData(filtersData []*alerts.FilterData) map[string][]string {
	result := make(map[string][]string, len(filtersData))
	for _, filter := range filtersData {
		field := filter.GetField()
		result[field] = flattenTracingFilter(filter.GetFilters())
	}
	return result
}

func flattenTracingFilter(filters []*alerts.Filters) []string {
	result := make([]string, 0)
	for _, f := range filters {
		values := f.GetValues()
		switch operator := f.GetOperator(); operator {
		case "notEquals", "contains", "startsWith", "endsWith":
			for i, val := range values {
				values[i] = fmt.Sprintf("filter:%s:%s", operator, val)
			}
		}
		result = append(result, values...)
	}
	return result
}

func flattenFlowAlert(flow *alerts.FlowCondition) *coralogixv1alpha1.Flow {
	stages := flattenFlowStages(flow.Stages)
	return &coralogixv1alpha1.Flow{
		Stages: stages,
	}
}

func flattenFlowStages(stages []*alerts.FlowStage) []coralogixv1alpha1.FlowStage {
	result := make([]coralogixv1alpha1.FlowStage, 0, len(stages))
	for _, s := range stages {
		stage := flattenFlowStage(s)
		result = append(result, stage)
	}
	return result
}

func flattenFlowStage(stage *alerts.FlowStage) coralogixv1alpha1.FlowStage {
	groups := flattenFlowStageGroups(stage.Groups)

	var timeFrame *coralogixv1alpha1.FlowStageTimeFrame
	if timeWindow := stage.GetTimeframe(); timeWindow != nil {
		timeFrame = convertMillisecondToTime(int(timeWindow.GetMs().GetValue()))
	}

	return coralogixv1alpha1.FlowStage{
		Groups:     groups,
		TimeWindow: timeFrame,
	}
}

func convertMillisecondToTime(timeMS int) *coralogixv1alpha1.FlowStageTimeFrame {
	if timeMS == 0 {
		return nil
	}

	msInHour := int(time.Hour.Milliseconds())
	msInMinute := int(time.Minute.Milliseconds())
	msInSecond := int(time.Second.Milliseconds())

	hours := timeMS / msInHour
	timeMS -= hours * msInHour

	minutes := timeMS / msInMinute
	timeMS -= minutes * msInMinute

	seconds := timeMS / msInSecond

	return &coralogixv1alpha1.FlowStageTimeFrame{
		Hours:   hours,
		Minutes: minutes,
		Seconds: seconds,
	}
}

func flattenFlowStageGroups(groups []*alerts.FlowGroup) []coralogixv1alpha1.FlowStageGroup {
	result := make([]coralogixv1alpha1.FlowStageGroup, 0, len(groups))
	for _, g := range groups {
		group := flattenFlowStageGroup(g)
		result = append(result, group)
	}
	return result
}

func flattenFlowStageGroup(group *alerts.FlowGroup) coralogixv1alpha1.FlowStageGroup {
	subAlerts := expandFlowSubgroupAlerts(group.GetAlerts())
	nextOp := alertProtoFlowOperatorToProtoFlowOperator[group.GetNextOp()]
	return coralogixv1alpha1.FlowStageGroup{
		InnerFlowAlerts: subAlerts,
		NextOperator:    nextOp,
	}
}

func expandFlowSubgroupAlerts(subgroup *alerts.FlowAlerts) coralogixv1alpha1.InnerFlowAlerts {
	return coralogixv1alpha1.InnerFlowAlerts{
		Operator: alertProtoFlowOperatorToProtoFlowOperator[subgroup.GetOp()],
		Alerts:   expandFlowInnerAlerts(subgroup.GetValues()),
	}
}

func expandFlowInnerAlerts(innerAlerts []*alerts.FlowAlert) []coralogixv1alpha1.InnerFlowAlert {
	result := make([]coralogixv1alpha1.InnerFlowAlert, 0, len(innerAlerts))
	for _, a := range innerAlerts {
		alert := expandFlowInnerAlert(a)
		result = append(result, alert)
	}
	return result
}

func expandFlowInnerAlert(alert *alerts.FlowAlert) coralogixv1alpha1.InnerFlowAlert {
	return coralogixv1alpha1.InnerFlowAlert{
		UserAlertId: alert.GetId().GetValue(),
		Not:         alert.GetNot().GetValue(),
	}
}

func flattenMetaLabels(labels []*alerts.MetaLabel) map[string]string {
	if len(labels) == 0 {
		return nil
	}

	result := make(map[string]string)
	for _, label := range labels {
		result[label.GetKey().GetValue()] = label.GetValue().GetValue()
	}
	return result
}

func flattenExpirationDate(expirationDate *alerts.Date) *coralogixv1alpha1.ExpirationDate {
	if expirationDate == nil {
		return nil
	}

	return &coralogixv1alpha1.ExpirationDate{
		Day:   expirationDate.Day,
		Month: expirationDate.Month,
		Year:  expirationDate.Year,
	}
}

func flattenScheduling(scheduling *alerts.AlertActiveWhen, timeZone coralogixv1alpha1.TimeZone) *coralogixv1alpha1.Scheduling {
	if scheduling == nil || len(scheduling.GetTimeframes()) == 0 {
		return nil
	}

	utc := coralogixv1alpha1.ExtractUTC(timeZone)
	timeframe := scheduling.GetTimeframes()[0]
	timeRange := timeframe.GetRange()
	activityStartGMT, activityEndGMT := timeRange.GetStart(), timeRange.GetEnd()
	daysOffset := getDaysOffsetFromGMT(activityStartGMT, utc)
	daysEnabled := flattenDaysOfWeek(timeframe.GetDaysOfWeek(), daysOffset)
	activityStartUTC := flattenTimeInDay(activityStartGMT, utc)
	activityEndUTC := flattenTimeInDay(activityEndGMT, utc)

	return &coralogixv1alpha1.Scheduling{
		TimeZone:    timeZone,
		DaysEnabled: daysEnabled,
		StartTime:   activityStartUTC,
		EndTime:     activityEndUTC,
	}
}

func getDaysOffsetFromGMT(activityStartGMT *alerts.Time, utc int32) int32 {
	daysOffset := int32(activityStartGMT.GetHours()+utc) / 24
	if daysOffset < 0 {
		daysOffset += 7
	}

	return daysOffset
}

func flattenTimeInDay(time *alerts.Time, utc int32) *coralogixv1alpha1.Time {
	hours := convertGmtToUtc(time.GetHours(), utc)
	hoursStr := toTwoDigitsFormat(hours)
	minStr := toTwoDigitsFormat(time.GetMinutes())
	result := coralogixv1alpha1.Time(fmt.Sprintf("%s:%s", hoursStr, minStr))
	return &result
}

func convertGmtToUtc(hours, utc int32) int32 {
	hours += utc
	if hours < 0 {
		hours += 24
	} else if hours >= 24 {
		hours -= 24
	}

	return hours
}

func toTwoDigitsFormat(digit int32) string {
	digitStr := fmt.Sprintf("%d", digit)
	if len(digitStr) == 1 {
		digitStr = "0" + digitStr
	}
	return digitStr
}

func flattenDaysOfWeek(daysOfWeek []alerts.DayOfWeek, daysOffset int32) []coralogixv1alpha1.Day {
	result := make([]coralogixv1alpha1.Day, 0, len(daysOfWeek))
	for _, d := range daysOfWeek {
		dayConvertedFromGmtToUtc := alerts.DayOfWeek((int32(d) + daysOffset) % 7)
		day := alertProtoDayToSchemaDay[dayConvertedFromGmtToUtc]
		result = append(result, day)
	}
	return result
}

func flattenNotificationGroups(ctx context.Context, log logr.Logger, notificationGroups []*alerts.AlertNotificationGroups) ([]coralogixv1alpha1.NotificationGroup, error) {
	result := make([]coralogixv1alpha1.NotificationGroup, 0, len(notificationGroups))
	webhooksIdsToNames, err := getWebhooksIdsToNames(ctx, log)
	if err != nil {
		return nil, fmt.Errorf("error on get webhooks ids to names - %w", err)
	}

	for _, ng := range notificationGroups {
		notificationGroup := flattenNotificationGroup(ng, webhooksIdsToNames)
		result = append(result, *notificationGroup)
	}

	return result, err
}

func getWebhooksIdsToNames(ctx context.Context, log logr.Logger) (map[uint32]string, error) {
	log.V(1).Info("get all webhooks")
	webhooks, err := coralogixv1alpha1.WebhooksClient.ListAllOutgoingWebhooks(ctx, &webhooks.ListAllOutgoingWebhooksRequest{})
	if err != nil {
		return nil, fmt.Errorf("error on get all webhooks - %w", err)
	}
	webhooksIdsToNames := make(map[uint32]string)
	for _, webhook := range webhooks.Deployed {
		webhooksIdsToNames[webhook.GetExternalId().GetValue()] = webhook.GetName().GetValue()
	}
	return webhooksIdsToNames, nil
}

func flattenNotificationGroup(notificationGroup *alerts.AlertNotificationGroups, webhooksIdsToNames map[uint32]string) *coralogixv1alpha1.NotificationGroup {
	return &coralogixv1alpha1.NotificationGroup{
		GroupByFields: utils.WrappedStringSliceToStringSlice(notificationGroup.GroupByFields),
		Notifications: flattenNotifications(notificationGroup.Notifications, webhooksIdsToNames),
	}
}

func flattenNotifications(notifications []*alerts.AlertNotification, webhooksIdsToNames map[uint32]string) []coralogixv1alpha1.Notification {
	result := make([]coralogixv1alpha1.Notification, 0, len(notifications))
	for _, notification := range notifications {
		flattenedNotification := flattenNotification(notification, webhooksIdsToNames)
		result = append(result, flattenedNotification)
	}
	return result
}

func flattenNotification(notification *alerts.AlertNotification, webhooksIdsToNames map[uint32]string) coralogixv1alpha1.Notification {
	notifyOn := alertProtoNotifyOn[notification.GetNotifyOn()]
	retriggeringPeriodMinutes := int32(notification.GetRetriggeringPeriodSeconds().GetValue()) / 60
	flattenedNotification := coralogixv1alpha1.Notification{
		NotifyOn:                  notifyOn,
		RetriggeringPeriodMinutes: retriggeringPeriodMinutes,
	}

	switch integration := notification.GetIntegrationType().(type) {
	case *alerts.AlertNotification_IntegrationId:
		webhookName, _ := webhooksIdsToNames[integration.IntegrationId.GetValue()]
		flattenedNotification.IntegrationName = pointer.String(webhookName)
	case *alerts.AlertNotification_Recipients:
		flattenedNotification.EmailRecipients = utils.WrappedStringSliceToStringSlice(integration.Recipients.Emails)
	}

	return flattenedNotification
}

func flattenShowInInsight(showInInsight *alerts.ShowInInsight) *coralogixv1alpha1.ShowInInsight {
	if showInInsight == nil {
		return nil
	}

	retriggeringPeriodMinutes := int32(showInInsight.GetRetriggeringPeriodSeconds().GetValue()) / 60
	notifyOn := alertProtoNotifyOn[showInInsight.GetNotifyOn()]

	return &coralogixv1alpha1.ShowInInsight{
		RetriggeringPeriodMinutes: retriggeringPeriodMinutes,
		NotifyOn:                  notifyOn,
	}
}

// SetupWithManager sets up the controller with the Manager.
func (r *AlertReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&coralogixv1alpha1.Alert{}).
		Complete(r)
}
