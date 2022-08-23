package claims

import (
	"time"

	"github.com/supertokens/supertokens-golang/supertokens"
)

func PrimitiveArrayClaim(key string, fetchValue FetchValueFunc) *TypePrimitiveArrayClaim {
	sessionClaim := SessionClaim(key, fetchValue)

	sessionClaim.AddToPayload_internal = func(payload map[string]interface{}, value interface{}, userContext supertokens.UserContext) map[string]interface{} {
		_, ok := value.([]interface{})
		assertCondition(ok, "value not an array")
		payload[sessionClaim.Key] = map[string]interface{}{
			"v": value,
			"t": time.Now().Unix(),
		}
		return payload
	}

	sessionClaim.RemoveFromPayloadByMerge_internal = func(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{} {
		payload[sessionClaim.Key] = nil
		return payload
	}

	sessionClaim.RemoveFromPayload = func(payload map[string]interface{}, userContext supertokens.UserContext) map[string]interface{} {
		delete(payload, sessionClaim.Key)
		return payload
	}

	sessionClaim.GetValueFromPayload = func(payload map[string]interface{}, userContext supertokens.UserContext) interface{} {
		if value, ok := payload[sessionClaim.Key].(map[string]interface{}); ok {
			return value["v"]
		}
		return nil
	}

	primitiveArrayClaim := &TypePrimitiveArrayClaim{
		TypeSessionClaim: sessionClaim,
	}

	primitiveArrayClaim.GetLastRefetchTime = func(payload map[string]interface{}, userContext supertokens.UserContext) *int64 {
		if value, ok := payload[sessionClaim.Key].(map[string]interface{}); ok {
			val := value["t"].(int64)
			return &val
		}
		return nil
	}

	primitiveArrayClaim.Validators = &PrimitiveArrayClaimValidators{
		Includes: func(val interface{}, maxAgeInSeconds *int64, id *string) *SessionClaimValidator {
			claimId := sessionClaim.Key + "-includes"
			if id != nil {
				claimId = *id
			}
			return &SessionClaimValidator{
				ID:    claimId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					claimVal, ok := sessionClaim.GetValueFromPayload(payload, userContext).(map[string]interface{})
					if !ok || claimVal == nil {
						return true
					}
					if maxAgeInSeconds != nil {
						return claimVal["t"].(int64) < time.Now().Unix()-*maxAgeInSeconds
					}
					return false
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal, claimValOk := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})
					assertCondition(claimValOk, "claim value not an array")

					if claimVal == nil {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "value does not exist",
								"expectedToInclude": val,
								"actualValue":       claimVal,
							},
						}
					}
					ageInSeconds := time.Now().Unix() - *primitiveArrayClaim.GetLastRefetchTime(payload, userContext)
					if maxAgeInSeconds != nil && ageInSeconds > *maxAgeInSeconds {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":         "expired",
								"ageInSeconds":    ageInSeconds,
								"maxAgeInSeconds": *maxAgeInSeconds,
							},
						}
					}
					if !includes(claimVal, val) {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "wrong value",
								"expectedToInclude": val,
								"actualValue":       claimVal,
							},
						}
					}
					return ClaimValidationResult{
						IsValid: true,
					}
				},
			}
		},
		Excludes: func(val interface{}, maxAgeInSeconds *int64, id *string) *SessionClaimValidator {
			claimId := sessionClaim.Key + "-excludes"
			if id != nil {
				claimId = *id
			}
			return &SessionClaimValidator{
				ID:    claimId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					val, ok := sessionClaim.GetValueFromPayload(payload, userContext).(map[string]interface{})
					if !ok || val == nil {
						return true
					}
					if maxAgeInSeconds != nil {
						return val["t"].(int64) < time.Now().Unix()-*maxAgeInSeconds
					}
					return false
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal, claimValOk := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})
					assertCondition(claimValOk, "claim value not an array")

					if claimVal == nil {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "value does not exist",
								"expectedToInclude": val,
								"actualValue":       claimVal,
							},
						}
					}
					ageInSeconds := time.Now().Unix() - *primitiveArrayClaim.GetLastRefetchTime(payload, userContext)
					if maxAgeInSeconds != nil && ageInSeconds > *maxAgeInSeconds {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":         "expired",
								"ageInSeconds":    ageInSeconds,
								"maxAgeInSeconds": *maxAgeInSeconds,
							},
						}
					}
					if includes(claimVal, val) {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "wrong value",
								"expectedToExclude": val,
								"actualValue":       claimVal,
							},
						}
					}
					return ClaimValidationResult{
						IsValid: true,
					}
				},
			}
		},
		IncludesAll: func(vals []interface{}, maxAgeInSeconds *int64, id *string) *SessionClaimValidator {
			claimId := sessionClaim.Key + "-includes-all"
			if id != nil {
				claimId = *id
			}
			return &SessionClaimValidator{
				ID:    claimId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					val, ok := sessionClaim.GetValueFromPayload(payload, userContext).(map[string]interface{})
					if !ok || val == nil {
						return true
					}
					if maxAgeInSeconds != nil {
						return val["t"].(int64) < time.Now().Unix()-*maxAgeInSeconds
					}
					return false
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal, claimValOk := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})
					assertCondition(claimValOk, "claim value not an array")

					if claimVal == nil {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "value does not exist",
								"expectedToInclude": vals,
								"actualValue":       claimVal,
							},
						}
					}
					ageInSeconds := time.Now().Unix() - *primitiveArrayClaim.GetLastRefetchTime(payload, userContext)
					if maxAgeInSeconds != nil && ageInSeconds > *maxAgeInSeconds {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":         "expired",
								"ageInSeconds":    ageInSeconds,
								"maxAgeInSeconds": *maxAgeInSeconds,
							},
						}
					}

					isValid := true
					valsMap := map[interface{}]bool{}
					for _, v := range vals {
						valsMap[v] = true
					}
					for _, v := range claimVal {
						if !valsMap[v] {
							isValid = false
							break
						}
					}

					if !isValid {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":           "wrong value",
								"expectedToInclude": vals,
								"actualValue":       claimVal,
							},
						}
					}
					return ClaimValidationResult{
						IsValid: true,
					}
				},
			}
		},
		ExcludesAll: func(vals []interface{}, maxAgeInSeconds *int64, id *string) *SessionClaimValidator {
			claimId := sessionClaim.Key + "-excludes-all"
			if id != nil {
				claimId = *id
			}
			return &SessionClaimValidator{
				ID:    claimId,
				Claim: sessionClaim,
				ShouldRefetch: func(payload map[string]interface{}, userContext supertokens.UserContext) bool {
					val, ok := sessionClaim.GetValueFromPayload(payload, userContext).(map[string]interface{})
					if !ok || val == nil {
						return true
					}
					if maxAgeInSeconds != nil {
						return val["t"].(int64) < time.Now().Unix()-*maxAgeInSeconds
					}
					return false
				},
				Validate: func(payload map[string]interface{}, userContext supertokens.UserContext) ClaimValidationResult {
					claimVal, claimValOk := sessionClaim.GetValueFromPayload(payload, userContext).([]interface{})
					assertCondition(claimValOk, "claim value not an array")

					if claimVal == nil {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":              "value does not exist",
								"expectedToNotInclude": vals,
								"actualValue":          claimVal,
							},
						}
					}
					ageInSeconds := time.Now().Unix() - *primitiveArrayClaim.GetLastRefetchTime(payload, userContext)
					if maxAgeInSeconds != nil && ageInSeconds > *maxAgeInSeconds {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":         "expired",
								"ageInSeconds":    ageInSeconds,
								"maxAgeInSeconds": *maxAgeInSeconds,
							},
						}
					}

					isValid := true
					valsMap := map[interface{}]bool{}
					for _, v := range vals {
						valsMap[v] = true
					}
					for _, v := range claimVal {
						if valsMap[v] {
							isValid = false
							break
						}
					}

					if !isValid {
						return ClaimValidationResult{
							IsValid: false,
							Reason: map[string]interface{}{
								"message":              "wrong value",
								"expectedToNotInclude": vals,
								"actualValue":          claimVal,
							},
						}
					}
					return ClaimValidationResult{
						IsValid: true,
					}
				},
			}
		},
	}

	return primitiveArrayClaim
}

type TypePrimitiveArrayClaim struct {
	*TypeSessionClaim
	GetLastRefetchTime func(payload map[string]interface{}, userContext supertokens.UserContext) *int64
	Validators         *PrimitiveArrayClaimValidators
}

type PrimitiveArrayClaimValidators struct {
	Includes    func(val interface{}, maxAgeInSeconds *int64, id *string) *SessionClaimValidator
	Excludes    func(val interface{}, maxAgeInSeconds *int64, id *string) *SessionClaimValidator
	IncludesAll func(vals []interface{}, maxAgeInSeconds *int64, id *string) *SessionClaimValidator
	ExcludesAll func(vals []interface{}, maxAgeInSeconds *int64, id *string) *SessionClaimValidator
}
