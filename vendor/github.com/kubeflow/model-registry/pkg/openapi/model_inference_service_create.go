/*
Model Registry REST API

REST API for Model Registry to create and manage ML model metadata

API version: v1alpha3
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package openapi

import (
	"encoding/json"
)

// checks if the InferenceServiceCreate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &InferenceServiceCreate{}

// InferenceServiceCreate An `InferenceService` entity in a `ServingEnvironment` represents a deployed `ModelVersion` from a `RegisteredModel` created by Model Serving.
type InferenceServiceCreate struct {
	// User provided custom properties which are not defined by its type.
	CustomProperties *map[string]MetadataValue `json:"customProperties,omitempty"`
	// An optional description about the resource.
	Description *string `json:"description,omitempty"`
	// The external id that come from the clients’ system. This field is optional. If set, it must be unique among all resources within a database instance.
	ExternalId *string `json:"externalId,omitempty"`
	// The client provided name of the artifact. This field is optional. If set, it must be unique among all the artifacts of the same artifact type within a database instance and cannot be changed once set.
	Name *string `json:"name,omitempty"`
	// ID of the `ModelVersion` to serve. If it's unspecified, then the latest `ModelVersion` by creation order will be served.
	ModelVersionId *string `json:"modelVersionId,omitempty"`
	// Model runtime.
	Runtime      *string                `json:"runtime,omitempty"`
	DesiredState *InferenceServiceState `json:"desiredState,omitempty"`
	// ID of the `RegisteredModel` to serve.
	RegisteredModelId string `json:"registeredModelId"`
	// ID of the parent `ServingEnvironment` for this `InferenceService` entity.
	ServingEnvironmentId string `json:"servingEnvironmentId"`
}

// NewInferenceServiceCreate instantiates a new InferenceServiceCreate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInferenceServiceCreate(registeredModelId string, servingEnvironmentId string) *InferenceServiceCreate {
	this := InferenceServiceCreate{}
	var desiredState InferenceServiceState = INFERENCESERVICESTATE_DEPLOYED
	this.DesiredState = &desiredState
	this.RegisteredModelId = registeredModelId
	this.ServingEnvironmentId = servingEnvironmentId
	return &this
}

// NewInferenceServiceCreateWithDefaults instantiates a new InferenceServiceCreate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInferenceServiceCreateWithDefaults() *InferenceServiceCreate {
	this := InferenceServiceCreate{}
	var desiredState InferenceServiceState = INFERENCESERVICESTATE_DEPLOYED
	this.DesiredState = &desiredState
	return &this
}

// GetCustomProperties returns the CustomProperties field value if set, zero value otherwise.
func (o *InferenceServiceCreate) GetCustomProperties() map[string]MetadataValue {
	if o == nil || IsNil(o.CustomProperties) {
		var ret map[string]MetadataValue
		return ret
	}
	return *o.CustomProperties
}

// GetCustomPropertiesOk returns a tuple with the CustomProperties field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InferenceServiceCreate) GetCustomPropertiesOk() (*map[string]MetadataValue, bool) {
	if o == nil || IsNil(o.CustomProperties) {
		return nil, false
	}
	return o.CustomProperties, true
}

// HasCustomProperties returns a boolean if a field has been set.
func (o *InferenceServiceCreate) HasCustomProperties() bool {
	if o != nil && !IsNil(o.CustomProperties) {
		return true
	}

	return false
}

// SetCustomProperties gets a reference to the given map[string]MetadataValue and assigns it to the CustomProperties field.
func (o *InferenceServiceCreate) SetCustomProperties(v map[string]MetadataValue) {
	o.CustomProperties = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *InferenceServiceCreate) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InferenceServiceCreate) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *InferenceServiceCreate) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *InferenceServiceCreate) SetDescription(v string) {
	o.Description = &v
}

// GetExternalId returns the ExternalId field value if set, zero value otherwise.
func (o *InferenceServiceCreate) GetExternalId() string {
	if o == nil || IsNil(o.ExternalId) {
		var ret string
		return ret
	}
	return *o.ExternalId
}

// GetExternalIdOk returns a tuple with the ExternalId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InferenceServiceCreate) GetExternalIdOk() (*string, bool) {
	if o == nil || IsNil(o.ExternalId) {
		return nil, false
	}
	return o.ExternalId, true
}

// HasExternalId returns a boolean if a field has been set.
func (o *InferenceServiceCreate) HasExternalId() bool {
	if o != nil && !IsNil(o.ExternalId) {
		return true
	}

	return false
}

// SetExternalId gets a reference to the given string and assigns it to the ExternalId field.
func (o *InferenceServiceCreate) SetExternalId(v string) {
	o.ExternalId = &v
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *InferenceServiceCreate) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InferenceServiceCreate) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *InferenceServiceCreate) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *InferenceServiceCreate) SetName(v string) {
	o.Name = &v
}

// GetModelVersionId returns the ModelVersionId field value if set, zero value otherwise.
func (o *InferenceServiceCreate) GetModelVersionId() string {
	if o == nil || IsNil(o.ModelVersionId) {
		var ret string
		return ret
	}
	return *o.ModelVersionId
}

// GetModelVersionIdOk returns a tuple with the ModelVersionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InferenceServiceCreate) GetModelVersionIdOk() (*string, bool) {
	if o == nil || IsNil(o.ModelVersionId) {
		return nil, false
	}
	return o.ModelVersionId, true
}

// HasModelVersionId returns a boolean if a field has been set.
func (o *InferenceServiceCreate) HasModelVersionId() bool {
	if o != nil && !IsNil(o.ModelVersionId) {
		return true
	}

	return false
}

// SetModelVersionId gets a reference to the given string and assigns it to the ModelVersionId field.
func (o *InferenceServiceCreate) SetModelVersionId(v string) {
	o.ModelVersionId = &v
}

// GetRuntime returns the Runtime field value if set, zero value otherwise.
func (o *InferenceServiceCreate) GetRuntime() string {
	if o == nil || IsNil(o.Runtime) {
		var ret string
		return ret
	}
	return *o.Runtime
}

// GetRuntimeOk returns a tuple with the Runtime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InferenceServiceCreate) GetRuntimeOk() (*string, bool) {
	if o == nil || IsNil(o.Runtime) {
		return nil, false
	}
	return o.Runtime, true
}

// HasRuntime returns a boolean if a field has been set.
func (o *InferenceServiceCreate) HasRuntime() bool {
	if o != nil && !IsNil(o.Runtime) {
		return true
	}

	return false
}

// SetRuntime gets a reference to the given string and assigns it to the Runtime field.
func (o *InferenceServiceCreate) SetRuntime(v string) {
	o.Runtime = &v
}

// GetDesiredState returns the DesiredState field value if set, zero value otherwise.
func (o *InferenceServiceCreate) GetDesiredState() InferenceServiceState {
	if o == nil || IsNil(o.DesiredState) {
		var ret InferenceServiceState
		return ret
	}
	return *o.DesiredState
}

// GetDesiredStateOk returns a tuple with the DesiredState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InferenceServiceCreate) GetDesiredStateOk() (*InferenceServiceState, bool) {
	if o == nil || IsNil(o.DesiredState) {
		return nil, false
	}
	return o.DesiredState, true
}

// HasDesiredState returns a boolean if a field has been set.
func (o *InferenceServiceCreate) HasDesiredState() bool {
	if o != nil && !IsNil(o.DesiredState) {
		return true
	}

	return false
}

// SetDesiredState gets a reference to the given InferenceServiceState and assigns it to the DesiredState field.
func (o *InferenceServiceCreate) SetDesiredState(v InferenceServiceState) {
	o.DesiredState = &v
}

// GetRegisteredModelId returns the RegisteredModelId field value
func (o *InferenceServiceCreate) GetRegisteredModelId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RegisteredModelId
}

// GetRegisteredModelIdOk returns a tuple with the RegisteredModelId field value
// and a boolean to check if the value has been set.
func (o *InferenceServiceCreate) GetRegisteredModelIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RegisteredModelId, true
}

// SetRegisteredModelId sets field value
func (o *InferenceServiceCreate) SetRegisteredModelId(v string) {
	o.RegisteredModelId = v
}

// GetServingEnvironmentId returns the ServingEnvironmentId field value
func (o *InferenceServiceCreate) GetServingEnvironmentId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ServingEnvironmentId
}

// GetServingEnvironmentIdOk returns a tuple with the ServingEnvironmentId field value
// and a boolean to check if the value has been set.
func (o *InferenceServiceCreate) GetServingEnvironmentIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ServingEnvironmentId, true
}

// SetServingEnvironmentId sets field value
func (o *InferenceServiceCreate) SetServingEnvironmentId(v string) {
	o.ServingEnvironmentId = v
}

func (o InferenceServiceCreate) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o InferenceServiceCreate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.CustomProperties) {
		toSerialize["customProperties"] = o.CustomProperties
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.ExternalId) {
		toSerialize["externalId"] = o.ExternalId
	}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.ModelVersionId) {
		toSerialize["modelVersionId"] = o.ModelVersionId
	}
	if !IsNil(o.Runtime) {
		toSerialize["runtime"] = o.Runtime
	}
	if !IsNil(o.DesiredState) {
		toSerialize["desiredState"] = o.DesiredState
	}
	toSerialize["registeredModelId"] = o.RegisteredModelId
	toSerialize["servingEnvironmentId"] = o.ServingEnvironmentId
	return toSerialize, nil
}

type NullableInferenceServiceCreate struct {
	value *InferenceServiceCreate
	isSet bool
}

func (v NullableInferenceServiceCreate) Get() *InferenceServiceCreate {
	return v.value
}

func (v *NullableInferenceServiceCreate) Set(val *InferenceServiceCreate) {
	v.value = val
	v.isSet = true
}

func (v NullableInferenceServiceCreate) IsSet() bool {
	return v.isSet
}

func (v *NullableInferenceServiceCreate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInferenceServiceCreate(val *InferenceServiceCreate) *NullableInferenceServiceCreate {
	return &NullableInferenceServiceCreate{value: val, isSet: true}
}

func (v NullableInferenceServiceCreate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInferenceServiceCreate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
