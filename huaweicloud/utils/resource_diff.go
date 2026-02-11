package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	awspolicy "github.com/jen20/awspolicyequivalence"
)

// ComposeAnySchemaDiffSuppressFunc allows parameters to determine multiple Diff methods.
// When any method (SchemaDiffSuppressFunc) returns true, this compose function will return true.
// It will only return false when all methods (SchemaDiffSuppressFunc) return false.
func ComposeAnySchemaDiffSuppressFunc(fs ...schema.SchemaDiffSuppressFunc) schema.SchemaDiffSuppressFunc {
	return func(k, o, n string, d *schema.ResourceData) bool {
		for _, f := range fs {
			if f(k, o, n, d) {
				return true
			}
		}
		return false
	}
}

func SuppressEquivalentAwsPolicyDiffs(k, old, new string, d *schema.ResourceData) bool {
	equivalent, err := awspolicy.PoliciesAreEquivalent(old, new)
	if err != nil {
		return false
	}

	return equivalent
}

// Suppress all changes
func SuppressDiffAll(k, old, new string, d *schema.ResourceData) bool {
	return true
}

// The SuppressCaseDiffs method compares the old and new values ​​of the current parameter to determine if their
// contents are consistent while ignoring the case format.
func SuppressCaseDiffs() schema.SchemaDiffSuppressFunc {
	return func(_, oldVal, newVal string, _ *schema.ResourceData) bool {
		return strings.EqualFold(oldVal, newVal)
	}
}

// Suppress changes if we get a computed min_disk_gb if value is unspecified (default 0)
func SuppressMinDisk(k, old, new string, d *schema.ResourceData) bool {
	return new == "0" || old == new
}

// Suppress changes if we get a base64 format or plaint text user_data
func SuppressUserData(k, old, new string, d *schema.ResourceData) bool {
	// user_data is in base64 format
	if HashAndHexEncode(old) == new {
		return true
	}

	// user_data is plaint text
	if plaint, err := base64.StdEncoding.DecodeString(old); err == nil {
		if HashAndHexEncode(string(plaint)) == new {
			return true
		}
	}

	return false
}

func SuppressTrimSpace(_, old, new string, _ *schema.ResourceData) bool {
	return strings.TrimSpace(old) == strings.TrimSpace(new)
}

func SuppressLBWhitelistDiffs(k, old, new string, d *schema.ResourceData) bool {
	if len(old) != len(new) {
		return false
	}
	old_array := strings.Split(old, ",")
	new_array := strings.Split(new, ",")
	sort.Strings(old_array)
	sort.Strings(new_array)

	return reflect.DeepEqual(old_array, new_array)
}

func SuppressSnatFiplistDiffs(k, old, new string, d *schema.ResourceData) bool {
	if len(old) != len(new) {
		return false
	}
	old_array := strings.Split(old, ",")
	new_array := strings.Split(new, ",")
	sort.Strings(old_array)
	sort.Strings(new_array)

	return reflect.DeepEqual(old_array, new_array)
}

// Suppress changes if we get a string with or without new line
func SuppressNewLineDiffs(k, old, new string, d *schema.ResourceData) bool {
	return strings.Trim(old, "\n") == strings.Trim(new, "\n")
}

func SuppressEquivilentTimeDiffs(k, old, new string, d *schema.ResourceData) bool {
	oldTime, err := time.Parse(time.RFC3339, old)
	if err != nil {
		return false
	}

	newTime, err := time.Parse(time.RFC3339, new)
	if err != nil {
		return false
	}

	return oldTime.Equal(newTime)
}

func SuppressVersionDiffs(k, old, new string, d *schema.ResourceData) bool {
	oldArray := regexp.MustCompile(`[\.\-]+`).Split(old, -1)
	newArray := regexp.MustCompile(`[\.\-]+`).Split(new, -1)
	if len(newArray) > len(oldArray) {
		return false
	}
	for i, v := range newArray {
		if v != oldArray[i] {
			return false
		}
	}
	return true
}

func CompareJsonTemplateAreEquivalent(tem1, tem2 string) (bool, error) {
	var obj1 interface{}
	err := json.Unmarshal([]byte(tem1), &obj1)
	if err != nil {
		return false, err
	}

	canonicalJson1, _ := json.Marshal(obj1)

	var obj2 interface{}
	err = json.Unmarshal([]byte(tem2), &obj2)
	if err != nil {
		return false, err
	}

	canonicalJson2, _ := json.Marshal(obj2)

	equal := bytes.Equal(canonicalJson1, canonicalJson2)
	if !equal {
		log.Printf("[DEBUG] Canonical template are not equal.\nFirst: %s\nSecond: %s\n",
			canonicalJson1, canonicalJson2)
	}
	return equal, nil
}

func SuppressStringSepratedByCommaDiffs(_, old, new string, _ *schema.ResourceData) bool {
	if len(old) != len(new) {
		return false
	}
	oldArray := strings.Split(old, ",")
	newArray := strings.Split(new, ",")
	sort.Strings(oldArray)
	sort.Strings(newArray)

	return reflect.DeepEqual(oldArray, newArray)
}

// ContainsAllKeyValues ​​checks whether object A (type map[string]interface{}) recursively contains all the keys and
// corresponding values ​​of object B (type map[string]interface{}).
// If the key-value pair in object B exists in object A and the values ​​are equal (recursively processing nested maps),
// return true; otherwise return false.
func ContainsAllKeyValues(objA, objB map[string]interface{}) bool {
	for key, bVal := range objB {
		aVal, exists := objA[key]
		if !exists {
			return false // A is missing the key of B.
		}

		// Check if the values ​​are both nested maps, if so, recursively compare.
		aMap, aIsMap := aVal.(map[string]interface{})
		bMap, bIsMap := bVal.(map[string]interface{})
		if aIsMap && bIsMap {
			if !ContainsAllKeyValues(aMap, bMap) {
				return false
			}
		} else {
			// Non-map types are compared directly via DeepEqual().
			if !reflect.DeepEqual(bVal, aVal) {
				return false
			}
		}
	}
	return true
}

// FindDecreaseKeys is a method that used to find out the key that objB is missing compared to objA.
// Will ignore the increase parts.
func FindDecreaseKeys(objA, objB map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, valA := range objA {
		if valB, exists := objB[key]; !exists {
			// If the key does not exist in objB, it's considered as a decrease key and is added directly to the result.
			result[key] = valA
		} else {
			// Check if the current values (valA and valB) are both type map for recursive processing.
			mapA, okA := valA.(map[string]interface{})
			mapB, okB := valB.(map[string]interface{})
			// If either valA or valB is not of type map, the subsequent recursive comparison is performed.
			if okA && okB {
				subResult := FindDecreaseKeys(mapA, mapB)
				if len(subResult) > 0 {
					result[key] = subResult
				}
			}
		}
	}
	return result
}

// SuppressObjectDiffs is a method that make the JSON string type parameter ignore the changes made on the console and
// only allow the local script to take effect.
func SuppressObjectDiffs() schema.SchemaDiffSuppressFunc {
	return func(paramKey, o, n string, d *schema.ResourceData) bool {
		if strings.HasSuffix(paramKey, ".%") || strings.HasSuffix(paramKey, ".#") {
			log.Printf("[DEBUG] The current change object is not of type object.")
			return false
		}
		return diffObjectParam(paramKey, o, n, d)
	}
}

// diffObjectParam is used to check whether the parameters of the current object or JSON object type have been modified
// other than those changed in the console.
// The following three scenarios will determine whether the parameter has changed (method return false):
//  1. The new value of the script adds some keys compared to the server return value (which must include keys that do
//     not exist in the value returned by the server).
//  2. The new value of the script modifies some (or all) key/value compared to the server return value.
//  3. The new value of the script removes some (or all) key/value compared to the old value of the script (the key can
//     be a nested structure).
//
// The following are examples of related scenarios:
//
// Service result:
//
//	{
//		"A": {
//			"Aa": "aa_aa",
//			"Ab": "aa_bb"
//		},
//		"B": "bb",
//		"C": "cc",
//		"D": "dd"
//	}
//
// Example 1 (Key 'D' add but the value is the same as the service result, so return true):
//
//	{					{
//		"B": "bb",			"B": "bb",
//		"C": "cc"	->		"C": "cc",
//	}						"D": "dd"
//						}
//
// Example 2 (New key 'D' addreturn false):
//
//	{					{
//		"B": "bb",			"B": "bb",
//		"C": "cc",	->		"C": "cc",
//	}						"E": "ee"
//						}
//
// Example 3 (The value of key 'C' changed, return false):
//
//	{					{
//		"B": "bb",			"B": "bb",
//		"C": "cc",	->		"C": "ccc"
//	}					}
//
// Example 4 (The value of key 'A.Aa' changed, return false):
//
//	{							{
//		"A": {						"A": {
//			"Aa": "aa_aa"				"Aa": "aa_aaa"
//		},					->		},
//		"B": "bb"					"B": "bb"
//	}							}
//
// Example 5 (Key 'D' removed, even it is exist in the service result, return false):
//
//	{					{
//		"B": "bb",			"B": "bb",
//		"C": "cc",	->		"C": "cc"
//		"D": "dd"		}
//	}
func diffObjectParam(paramKey, _, _ string, d *schema.ResourceData) bool {
	var (
		consoleVal, newScriptVal, originVal map[string]interface{}

		originParamKey           = fmt.Sprintf("%s_origin", paramKey)
		oldParamVal, newParamVal = d.GetChange(paramKey)
	)

	// After refresh phase, the value from the service side will be stored in the tfstate, and as old value in the
	// next d.GetChange() method returns.
	consoleVal = TryMapValueAnalysis(oldParamVal)
	newScriptVal = TryMapValueAnalysis(newParamVal)
	// The script value of the last update (if it has) is used as a reference for the historical value of this
	// change.
	originVal = TryMapValueAnalysis(d.Get(originParamKey))

	return ContainsAllKeyValues(consoleVal, newScriptVal) && len(FindDecreaseKeys(originVal, newScriptVal)) < 1
}

func SuppressMapDiffs() schema.SchemaDiffSuppressFunc {
	return func(paramKey, o, n string, d *schema.ResourceData) bool {
		log.Printf("[DEBUG][SuppressMapDiffs] Called with paramKey='%s', old='%s', new='%s'", paramKey, o, n)

		// Ignore length change judgment, because this method will judge each changed key one by one
		if strings.HasSuffix(paramKey, ".%") {
			log.Printf("[DEBUG][SuppressMapDiffs] Ignoring length change for '%s'", paramKey)
			return true
		}

		// Handle the case where the entire map is being changed
		if !strings.Contains(paramKey, ".") {
			return suppressEntireMapDiff(paramKey, d)
		}

		// Handle single key changes
		return suppressSingleKeyDiff(paramKey, d)
	}
}

// suppressEntireMapDiff handles changes to the entire map
func suppressEntireMapDiff(paramKey string, d *schema.ResourceData) bool {
	originMapCategory := fmt.Sprintf("%s_origin", paramKey)
	log.Printf("[DEBUG][EntireMapDiff] Handling entire map change for '%s', origin map category: '%s'",
		paramKey, originMapCategory)

	originMapVal := d.Get(originMapCategory)
	if originMapVal == nil {
		log.Printf("[DEBUG][EntireMapDiff] Origin map '%s' is nil, suppressing diff for entire map '%s'",
			originMapCategory, paramKey)
		return true
	}

	originMap := TryMapValueAnalysis(originMapVal)
	if len(originMap) == 0 {
		log.Printf("[DEBUG][EntireMapDiff] Origin map '%s' is empty, suppressing diff for entire map '%s'",
			originMapCategory, paramKey)
		return true
	}

	// For entire map changes, we need to check if all keys in the new value exist in origin
	// This is a simplified approach - in practice, you might want more sophisticated comparison
	log.Printf("[DEBUG][EntireMapDiff] Entire map '%s' change detected, checking against origin", paramKey)
	return false // For now, report the change and let individual key suppression handle it
}

// suppressSingleKeyDiff handles changes to a single key
func suppressSingleKeyDiff(paramKey string, d *schema.ResourceData) bool {
	categories := strings.Split(paramKey, ".")
	mapCategory := strings.Join(categories[:len(categories)-1], ".")
	originMapCategory := fmt.Sprintf("%s_origin", mapCategory)
	keyName := categories[len(categories)-1]

	log.Printf("[DEBUG][SingleKeyDiff] Processing key '%s', mapCategory='%s', originMapCategory='%s', keyName='%s'",
		paramKey, mapCategory, originMapCategory, keyName)

	// Get origin map (last local configuration)
	originMapVal := d.Get(originMapCategory)
	originMap := TryMapValueAnalysis(originMapVal)
	log.Printf("[DEBUG][SingleKeyDiff] Origin map '%s' content: %+v", originMapCategory, originMap)

	// Get current configuration map
	currentMapVal := GetNestedObjectFromRawConfig(d.GetRawConfig(), mapCategory)
	if currentMapVal == nil {
		log.Printf("[DEBUG][SingleKeyDiff] Current map '%s' is nil, suppressing diff for key '%s'",
			mapCategory, keyName)
		return true
	}

	currentMap := TryMapValueAnalysis(currentMapVal)
	log.Printf("[DEBUG][SingleKeyDiff] Current map '%s' content: %+v", mapCategory, currentMap)

	// Check if the key exists in current configuration
	existsInCurrent := keyExists(currentMap, keyName)
	existsInOrigin := keyExists(originMap, keyName)
	isOriginEmpty := originMapVal == nil || len(originMap) == 0

	// Determine suppression based on key existence
	return determineSuppression(existsInCurrent, existsInOrigin, isOriginEmpty, keyName)
}

// keyExists checks if a key exists in the map
func keyExists(m map[string]interface{}, key string) bool {
	_, exists := m[key]
	return exists
}

// determineSuppression determines whether to suppress the diff based on key existence
func determineSuppression(existsInCurrent, existsInOrigin, isOriginEmpty bool, keyName string) bool {
	if existsInCurrent {
		// The key exists in current configuration
		if isOriginEmpty {
			// Origin is empty or nil, report the change (locally added)
			log.Printf("[DEBUG] The key '%s' found in current config but origin is empty", keyName)
			return false
		}

		if existsInOrigin {
			// The key exists in both current config and origin, report this change
			log.Printf("[DEBUG] The key '%s' found in both current config and origin", keyName)
			return false
		}

		// The key exists in current config but not in origin (locally added)
		log.Printf("[DEBUG] The key '%s' found in current config but not in origin", keyName)
		return false
	}

	// The key does not exist in current configuration
	if isOriginEmpty {
		// Origin is empty or nil, suppress the diff (remotely added)
		log.Printf("[DEBUG] The key '%s' not found in current config and origin is empty, suppressing diff",
			keyName)
		return true
	}

	if existsInOrigin {
		// The key exists in origin but not in current config (locally removed)
		log.Printf("[DEBUG] The key '%s' found in origin but not in current config (locally removed)",
			keyName)
		return false
	}

	// The key does not exist in either current config or origin (remotely added)
	log.Printf("[DEBUG] The key '%s' not found in either current config or origin (remotely added), suppressing diff",
		keyName)
	return true
}

// TakeObjectsDifferent is a method that used to recursively get the complement of object A (objA) compared to
// object B (objB) (including nested differences).
// In Math, it means A-B, also A\B or {x | x∈A，x∉B}.
func TakeObjectsDifferent(objA, objB map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for key, valueA := range objA {
		valueB, exists := objB[key]

		// The key in objA does not exist in objB
		if !exists {
			result[key] = valueA
			continue
		}

		// Try recursively processing nested map.
		if subMapA, okA := valueA.(map[string]interface{}); okA {
			if subMapB, okB := valueB.(map[string]interface{}); okB {
				// Both sides are maps, recursive comparison.
				subDiff := TakeObjectsDifferent(subMapA, subMapB)
				if len(subDiff) > 0 {
					result[key] = subDiff
				}
			} else {
				// The value of objA is a map but the value of objB is not (type inconsistency).
				result[key] = valueA
			}
			continue
		}

		// Handling non-map types or inconsistent types.
		if !reflect.DeepEqual(valueA, valueB) {
			result[key] = valueA
		}
	}

	return result
}

// SuppressStrSliceDiffs is a method that makes the string slice type parameter ignore the changes made on the console and
// only allow the local script to take effect. It identifies elements that are decreased compared to origin and
// elements that are newly added remotely.
func SuppressStrSliceDiffs() schema.SchemaDiffSuppressFunc {
	return func(paramKey, o, n string, d *schema.ResourceData) bool {
		log.Printf("[DEBUG][SuppressStrSliceDiffs] Called with paramKey='%s', oldVal='%s', newVal='%s'", paramKey, o, n)

		// Handle TypeSet length field
		if strings.HasSuffix(paramKey, ".#") {
			log.Printf("[DEBUG][SuppressStrSliceDiffs] Processing TypeSet length field: %s", paramKey)
			return diffStrSliceLength(paramKey, o, n, d)
		}

		// Handle TypeSet element fields (e.g., {set_param_key}.1234567890)
		if strings.Contains(paramKey, ".") && !strings.HasSuffix(paramKey, ".%") {
			log.Printf("[DEBUG][SuppressStrSliceDiffs] Processing TypeSet element field: %s", paramKey)
			return diffStrSliceElement(paramKey, o, n, d)
		}

		if strings.HasSuffix(paramKey, ".%") {
			log.Printf("[DEBUG][SuppressStrSliceDiffs] The current change object is not of type slice.")
			return false
		}

		log.Printf("[DEBUG][SuppressStrSliceDiffs] Processing main field: %s", paramKey)
		result := diffStrSliceParam(paramKey, o, n, d)
		log.Printf("[DEBUG][SuppressStrSliceDiffs] Final result: %v", result)
		return result
	}
}

// diffStrSliceLength handles the length field of TypeList or TypeSet
func diffStrSliceLength(paramKey, oldVal, newVal string, d *schema.ResourceData) bool {
	baseField := strings.TrimSuffix(paramKey, ".#")

	// Get the origin value
	originParamKey := fmt.Sprintf("%s_origin", baseField)
	originVal := d.Get(originParamKey)

	// Get current values
	oldCount, _ := strconv.Atoi(oldVal)
	newCount, _ := strconv.Atoi(newVal)

	// If origin is empty or nil, this is the first time setting the value
	if originVal == nil {
		return false
	}

	// Check if origin is effectively empty
	var originCount int
	var isEmpty bool
	switch v := originVal.(type) {
	case []interface{}:
		originCount = len(v)
		isEmpty = len(v) == 0
	case *schema.Set:
		originCount = v.Len()
		isEmpty = v.Len() == 0
	default:
		originCount = 0
		isEmpty = true
	}

	// If origin is empty, check if this is a remote-only change that should be suppressed
	if isEmpty {
		// Get current remote state to check if this is a remote-only change
		currentVal := d.Get(baseField)
		if currentVal != nil {
			var currentCount int
			switch v := currentVal.(type) {
			case []interface{}:
				currentCount = len(v)
			case *schema.Set:
				currentCount = v.Len()
			default:
				currentCount = 0
			}

			// If new count is less than current count, this might be a remote removal
			// that should be suppressed (unless it's a local removal)
			if newCount < currentCount {
				return true
			}
		}

		return false
	}

	// Check if there are actual changes that require length difference to be shown
	hasLocalAdditions := newCount > oldCount
	hasLocalRemovals := newCount < originCount

	// If there are actual local changes, don't suppress length difference
	if hasLocalAdditions || hasLocalRemovals {
		log.Printf("[DEBUG][diffStrSliceLength] Is local additions happened? %v", hasLocalAdditions)
		log.Printf("[DEBUG][diffStrSliceLength] Is local removals happened? %v", hasLocalRemovals)
		return false
	}

	// If no actual changes, suppress length difference (e.g., remote-only additions/removals)
	return true
}

// diffStrSliceElement handles individual slice elements. And for the TypeSet, there are indexes of each element,
// so we need to handle them separately.
func diffStrSliceElement(paramKey, oldVal, newVal string, d *schema.ResourceData) bool {
	parts := strings.Split(paramKey, ".")
	if len(parts) != 2 {
		log.Printf("[DEBUG][diffStrSliceElement] Invalid paramKey format: %s", paramKey)
		return false
	}
	baseField := parts[0]

	originParamKey := fmt.Sprintf("%s_origin", baseField)
	originVal := d.Get(originParamKey)

	log.Printf("[DEBUG][diffStrSliceElement] baseField='%s', oldVal='%s', newVal='%s', originVal=%v",
		baseField, oldVal, newVal, originVal)

	// Handle element removal case
	if newVal == "" {
		return handleElementRemoval(oldVal, originVal, baseField, d)
	}

	// Handle element addition/modification case
	return handleElementAddition(oldVal, newVal, originVal, baseField, d)
}

// handleElementRemoval handles the case when an element is being removed
func handleElementRemoval(oldVal string, originVal interface{}, baseField string, d *schema.ResourceData) bool {
	log.Printf("[DEBUG][handleElementRemoval] Element '%s' is being removed, checking if should suppress diff", oldVal)

	// Check if this element was in origin
	if isElementInOrigin(oldVal, originVal) {
		log.Printf("[DEBUG][handleElementRemoval] Element '%s' was in origin, NOT suppressing diff (allow removal)", oldVal)
		return false // NOT suppressing - allow removal of origin elements
	}

	// If origin is empty or nil, check if this element exists in remote state
	if checkElementInRemoteState(baseField, oldVal, d) {
		log.Printf("[DEBUG][handleElementRemoval] Element '%s' exists in remote state but not in origin, suppressing diff (ignore remote removal)",
			oldVal)
		return true // Suppress diff - ignore removal of remote-only elements
	}

	// If element was not in origin and not in remote state, suppress the diff
	log.Printf("[DEBUG][handleElementRemoval] Element '%s' was not in origin or remote state, suppressing diff",
		oldVal)
	return true
}

// handleElementAddition handles the case when an element is being added or modified
func handleElementAddition(oldVal, newVal string, originVal interface{}, baseField string, d *schema.ResourceData) bool {
	// If origin is nil or empty, this is the first time setting the value
	if isOriginEmpty(originVal) {
		return handleFirstTimeSetting(oldVal, newVal, baseField, d)
	}

	// Check if this element is in origin
	if isElementInOrigin(newVal, originVal) {
		// If the element is unchanged (oldVal == newVal), don't suppress diff
		// This ensures Terraform knows the config value still exists
		if oldVal == newVal {
			log.Printf("[DEBUG][handleElementAddition] Element '%s' unchanged and in origin, not suppressing diff to preserve config value", newVal)
			return false
		}
		log.Printf("[DEBUG][handleElementAddition] Element '%s' was in origin, suppressing diff", newVal)
		return true
	}

	// If element was not in origin, don't suppress (this is a local addition)
	log.Printf("[DEBUG][handleElementAddition] Element '%s' was not in origin, not suppressing diff (local addition)", newVal)
	return false
}

// handleFirstTimeSetting handles the case when origin is empty or nil
func handleFirstTimeSetting(oldVal, newVal, baseField string, d *schema.ResourceData) bool {
	// If oldVal is empty, this is a CREATE scenario - use main logic
	if oldVal == "" {
		log.Printf("[DEBUG][handleFirstTimeSetting] This is a CREATE scenario (oldVal=''), using main diffStrSliceParam logic")
		return false // Let the main logic handle this
	}
	// If oldVal is not empty, this is an UPDATE scenario - check remote state
	log.Printf("[DEBUG][handleFirstTimeSetting] This is an UPDATE scenario (oldVal='%s'), checking if new value exists in remote state", oldVal)
	return checkElementInRemoteState(baseField, newVal, d)
}

// isOriginEmpty checks if origin value is effectively empty
func isOriginEmpty(originVal interface{}) bool {
	if originVal == nil {
		return true
	}

	switch v := originVal.(type) {
	case []interface{}:
		return len(v) == 0
	case *schema.Set:
		return v.Len() == 0
	default:
		log.Printf("[DEBUG][isOriginEmpty] Unexpected originVal type: %T", originVal)
		return true
	}
}

// isElementInOrigin checks if an element exists in origin value
func isElementInOrigin(element string, originVal interface{}) bool {
	if originVal == nil {
		return false
	}

	switch v := originVal.(type) {
	case []interface{}:
		for _, item := range v {
			if str, ok := item.(string); ok && str == element {
				return true
			}
		}
	case *schema.Set:
		return v.Contains(element)
	}

	return false
}

// diffStrSliceParam is used to check whether the parameters of the current string slice type have been modified
// other than those changed in the console.
// The following scenarios will determine whether the parameter has changed (method return false):
//  1. The new value of the script adds new elements compared to the console value (locally added elements).
//  2. The new value of the script has elements decreased compared to the origin value (locally removed elements).
//
// The following scenarios will suppress the diff (method return true):
//  1. The new value of the script is a subset of the console value AND
//     the new value of the script has no elements decreased compared to the origin value.
//
// Examples:
//
// Origin value: ["a", "b", "c"]
// Console value: ["a", "b", "c", "d"] (remotely added "d")
//
// Example 1 (Return false - locally added new element):
//
//	Script value: ["a", "b", "c", "e"] -> Contains "e" not in console (locally added)
//
// Example 2 (Return false - locally removed element):
//
//	Script value: ["a", "b"] -> Removed "c" from origin (locally removed)
//
// Example 3 (Return true - subset of console and no decrease from origin):
//
//	Script value: ["a", "b", "c"] -> Subset of console, same as origin
//
// Example 4 (Return true - subset of console and no decrease from origin):
//
//	Script value: ["a", "b"] -> Subset of console, subset of origin (allowed decrease)
//
// Example 5 (Return false - locally removed element, even if console has new elements):
//
//	Origin: ["a", "b", "c"]
//	Console: ["a", "b", "c", "d", "e"] (remotely added "d", "e")
//	Script: ["a", "b"] -> Removed "c" from origin (locally removed)
//
// Example 6 (Return true - ignore remotely added elements, no local changes):
//
//	Origin: ["a", "b", "c"]
//	Console: ["a", "b", "c", "d", "e"] (remotely added "d", "e")
//	Script: ["a", "b", "c"] -> Same as origin, ignore remote additions
func diffStrSliceParam(paramKey, oldVal, newVal string, d *schema.ResourceData) bool {
	var (
		originSlice, consoleSlice, newScriptSlice []string
		originParamKey                            = fmt.Sprintf("%s_origin", paramKey)
	)

	// Get the origin value (last local configuration) - this is the key fix
	originVal := d.Get(originParamKey)
	if originVal != nil {
		// Handle different types that originVal might be
		switch v := originVal.(type) {
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok && str != "" {
					originSlice = append(originSlice, str)
				}
			}
		case *schema.Set:
			for _, item := range v.List() {
				if str, ok := item.(string); ok && str != "" {
					originSlice = append(originSlice, str)
				}
			}
		case string:
			if v != "" {
				originSlice = strings.Split(v, ",")
				originSlice = removeEmptyStrings(originSlice)
			}
		default:
			log.Printf("[DEBUG][diffStrSliceParam] Unexpected originVal type: %T", originVal)
		}
	}

	// If origin is empty, this is the first time setting the value
	// In this case, we should allow the change (not suppress diff)
	if len(originSlice) == 0 {
		log.Printf("[DEBUG][diffStrSliceParam] Origin is empty, allowing change (first time setting)")
		return false
	}

	// Parse the old and new values from GetChange
	// oldVal and newVal are already strings from Terraform's diff suppression
	// They represent the serialized form of the lists
	if oldVal != "" {
		consoleSlice = strings.Split(oldVal, ",")
		consoleSlice = removeEmptyStrings(consoleSlice)
	}
	if newVal != "" {
		newScriptSlice = strings.Split(newVal, ",")
		newScriptSlice = removeEmptyStrings(newScriptSlice)
	}

	log.Printf("[DEBUG][diffStrSliceParam] paramKey='%s', originSlice=%v, consoleSlice=%v, newScriptSlice=%v",
		paramKey, originSlice, consoleSlice, newScriptSlice)

	// Check if only care about elements that are in new script but NOT in console (locally added)
	// This means we ignore elements that are in console but NOT in new script (remotely added)
	localAdditions := FindStrSliceElementsNotInAnother(newScriptSlice, consoleSlice)
	if len(localAdditions) > 0 {
		log.Printf("[DEBUG][diffStrSliceParam] New script contains elements not in console (locally added): %v, not suppressing diff",
			localAdditions)
		return false
	}

	// Check if new script has elements decreased compared to origin (locally removed)
	// These are elements that will be deleted from remote
	log.Printf("[DEBUG][diffStrSliceParam] comparing newScriptSlice=%v with originSlice=%v", newScriptSlice, originSlice)
	localRemovals := FindStrSliceElementsNotInAnother(originSlice, newScriptSlice)
	if len(localRemovals) > 0 {
		log.Printf("[DEBUG][diffStrSliceParam] New script has elements decreased compared to origin (locally removed), not suppressing diff")
		return false
	}

	// Both conditions are met, suppress the diff
	log.Printf("[DEBUG][diffStrSliceParam] No local additions and no local removals, suppressing diff")
	return true
}

// removeEmptyStrings removes empty strings from a slice
func removeEmptyStrings(slice []string) []string {
	result := make([]string, 0, len(slice))
	for _, s := range slice {
		if strings.TrimSpace(s) != "" {
			result = append(result, strings.TrimSpace(s))
		}
	}
	return result
}

// checkElementInRemoteState checks if an element exists in remote state
func checkElementInRemoteState(baseField, elementValue string, d *schema.ResourceData) bool {
	// Get the current remote state value
	currentVal := d.Get(baseField)
	if currentVal != nil {
		switch v := currentVal.(type) {
		case []interface{}:
			for _, item := range v {
				if str, ok := item.(string); ok && str == elementValue {
					log.Printf("[DEBUG][checkElementInRemoteState] Element '%s' already exists in remote state, suppressing diff",
						elementValue)
					return true
				}
			}
		case *schema.Set:
			if v.Contains(elementValue) {
				log.Printf("[DEBUG][checkElementInRemoteState] Element '%s' already exists in remote state, suppressing diff",
					elementValue)
				return true
			}
		}
	}

	log.Printf("[DEBUG][checkElementInRemoteState] Element '%s' does not exist in remote state, allowing change",
		elementValue)
	return false
}

// SuppressObjectSliceDiffs is a method that makes the object slice type parameter ignore the changes made on the console and
// only allow the local script to take effect. It identifies elements that are decreased compared to origin and
// elements that are newly added remotely.
func SuppressObjectSliceDiffs() schema.SchemaDiffSuppressFunc {
	return func(paramKey, o, n string, d *schema.ResourceData) bool {
		log.Printf("[DEBUG][SuppressObjectSliceDiffs] Called with paramKey='%s', oldVal='%s', newVal='%s'", paramKey, o, n)

		// Handle TypeSet length field
		if strings.HasSuffix(paramKey, ".#") {
			log.Printf("[DEBUG][SuppressObjectSliceDiffs] Processing TypeSet length field: %s", paramKey)
			return diffObjectSliceLength(paramKey, o, n, d)
		}

		// Handle TypeSet element fields (e.g., {set_param_key}.1234567890.type, {set_param_key}.1234567890.id)
		if strings.Contains(paramKey, ".") && !strings.HasSuffix(paramKey, ".%") {
			log.Printf("[DEBUG][SuppressObjectSliceDiffs] Processing TypeSet element field: %s", paramKey)
			return diffObjectSliceElement(paramKey, o, n, d)
		}

		if strings.HasSuffix(paramKey, ".%") {
			log.Printf("[DEBUG][SuppressObjectSliceDiffs] The current change object is not of type slice.")
			return false
		}

		log.Printf("[DEBUG][SuppressObjectSliceDiffs] Processing main field: %s", paramKey)
		result := diffObjectSliceParam(paramKey, d)
		log.Printf("[DEBUG][SuppressObjectSliceDiffs] Final result: %v", result)
		return result
	}
}

// diffObjectSliceLength handles the length field of TypeList or TypeSet for object slices
func diffObjectSliceLength(paramKey, oldVal, newVal string, d *schema.ResourceData) bool {
	baseField := strings.TrimSuffix(paramKey, ".#")

	// Get the origin value
	originParamKey := fmt.Sprintf("%s_origin", baseField)
	originVal := d.Get(originParamKey)

	// Get current values
	oldCount, _ := strconv.Atoi(oldVal)
	newCount, _ := strconv.Atoi(newVal)

	// If counts are the same, suppress diff
	if oldCount == newCount {
		log.Printf("[DEBUG][diffObjectSliceLength] Old count (%d) equals new count (%d), suppressing diff", oldCount, newCount)
		return true
	}

	// Get old and new values from GetChange
	oldParamVal, newParamVal := d.GetChange(baseField)
	tfstateSlice := convertToObjectSlice(oldParamVal)   // tfstate (remote state)
	rawConfigSlice := convertToObjectSlice(newParamVal) // RawConfig (script config)
	originSlice := convertToObjectSlice(originVal)      // origin

	log.Printf("[DEBUG][diffObjectSliceLength] paramKey='%s', oldCount=%d, newCount=%d, originSlice length=%d, tfstateSlice length=%d, "+
		"rawConfigSlice length=%d",
		paramKey, oldCount, newCount, len(originSlice), len(tfstateSlice), len(rawConfigSlice))

	// Detailed logging for debugging
	log.Printf("[DEBUG][diffObjectSliceLength] 1. RawConfig: %v", formatObjectSliceForLog(rawConfigSlice))
	log.Printf("[DEBUG][diffObjectSliceLength] 2. tfstate: %v", formatObjectSliceForLog(tfstateSlice))
	log.Printf("[DEBUG][diffObjectSliceLength] 3. origin: %v", formatObjectSliceForLog(originSlice))

	// Check 1: (RawConfig - tfstate) ∩ (RawConfig - origin)
	// If not empty, it means there are local additions
	localAdditions := calculateLocalAdditions(originVal, oldParamVal, newParamVal)
	rawConfigMinusTfstate := FindObjectSliceElementsNotInAnother(rawConfigSlice, tfstateSlice)
	rawConfigMinusOrigin := FindObjectSliceElementsNotInAnother(rawConfigSlice, originSlice)

	log.Printf("[DEBUG][diffObjectSliceLength] 4. (RawConfig - tfstate): %v", formatObjectSliceForLog(rawConfigMinusTfstate))
	log.Printf("[DEBUG][diffObjectSliceLength] 5. (RawConfig - origin): %v", formatObjectSliceForLog(rawConfigMinusOrigin))
	log.Printf("[DEBUG][diffObjectSliceLength] 6. (RawConfig - tfstate) ∩ (RawConfig - origin): %v", formatObjectSliceForLog(localAdditions))

	if len(localAdditions) > 0 {
		log.Printf("[DEBUG][diffObjectSliceLength] Found local additions (RawConfig - tfstate) ∩ (RawConfig - origin): %d, NOT suppressing diff",
			len(localAdditions))
		return false
	}

	// Check 2: (tfstate - RawConfig) ∩ (origin - RawConfig)
	// If not empty, it means there are local removals
	localRemovals := calculateLocalRemovals(originVal, oldParamVal, newParamVal)
	tfstateMinusRawConfig := FindObjectSliceElementsNotInAnother(tfstateSlice, rawConfigSlice)
	originMinusRawConfig := FindObjectSliceElementsNotInAnother(originSlice, rawConfigSlice)

	log.Printf("[DEBUG][diffObjectSliceLength] 7. (tfstate - RawConfig): %v", formatObjectSliceForLog(tfstateMinusRawConfig))
	log.Printf("[DEBUG][diffObjectSliceLength] 8. (origin - RawConfig): %v", formatObjectSliceForLog(originMinusRawConfig))
	log.Printf("[DEBUG][diffObjectSliceLength] 9. (tfstate - RawConfig) ∩ (origin - RawConfig): %v", formatObjectSliceForLog(localRemovals))

	// Check if all removed elements are local removals
	// If there are remote additions (elements in tfstate - RawConfig but not in localRemovals),
	// we should suppress the length diff to allow element-level suppression
	remoteAdditions := make([]map[string]interface{}, 0)
	for _, elem := range tfstateMinusRawConfig {
		if !ObjectSliceContains(localRemovals, elem) {
			remoteAdditions = append(remoteAdditions, elem)
		}
	}

	log.Printf("[DEBUG][diffObjectSliceLength] Remote additions (elements in tfstate - RawConfig but not in localRemovals): %v",
		formatObjectSliceForLog(remoteAdditions))

	if len(localRemovals) > 0 {
		if len(remoteAdditions) > 0 {
			// There are both local removals and remote additions
			// We need to NOT suppress the length diff to show the deletion operation
			// The adjusted old count (excluding remote additions) should be used conceptually:
			// adjustedOldCount = oldCount - len(remoteAdditions) = 4 - 1 = 3
			// This means we want to show "3 -> 2" conceptually, but Terraform will show "4 -> 2"
			// However, we can suppress remote additions at element level, so only local removals are shown
			// Note: We cannot modify oldVal/newVal in diff suppression functions, but we can control
			// which elements are shown/hidden at element level
			adjustedOldCount := oldCount - len(remoteAdditions)
			log.Printf("[DEBUG][diffObjectSliceLength] Found local removals (%d) and remote additions (%d), "+
				"adjusted old count: %d -> %d (conceptually %d -> %d)",
				len(localRemovals), len(remoteAdditions), oldCount, newCount, adjustedOldCount, newCount)
			log.Printf("[DEBUG][diffObjectSliceLength] NOT suppressing length diff to show deletion operation, " +
				"remote additions will be suppressed at element level")
			return false
		}
		// All removed elements are local removals, don't suppress
		log.Printf("[DEBUG][diffObjectSliceLength] Found local removals (tfstate - RawConfig) ∩ (origin - RawConfig): %d, NOT suppressing diff",
			len(localRemovals))
		return false
	}

	// Both checks are empty, suppress diff
	log.Printf("[DEBUG][diffObjectSliceLength] No local additions or removals, suppressing diff")
	return true
}

// findTargetObjectFromState finds the target object from state attributes by objectHash
func findTargetObjectFromState(d *schema.ResourceData, baseField, objectHash string, tfstateSlice []map[string]interface{}) map[string]interface{} {
	if d.State() == nil || d.State().Attributes == nil {
		return nil
	}

	tempObject := make(map[string]interface{})
	for key, val := range d.State().Attributes {
		if strings.HasPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash)) {
			fieldName := strings.TrimPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash))
			tempObject[fieldName] = val
		}
	}

	if len(tempObject) == 0 {
		return nil
	}

	// Found object in state attributes, try to find the full object in tfstateSlice
	for _, tfstateObjMap := range tfstateSlice {
		matches := true
		for key, val := range tempObject {
			if tfstateVal, ok := tfstateObjMap[key]; !ok || !reflect.DeepEqual(val, tfstateVal) {
				matches = false
				break
			}
		}
		if matches {
			log.Printf("[DEBUG][findTargetObjectFromState] Found full object from state attributes matching objectHash '%s': %v",
				objectHash, tfstateObjMap)
			return tfstateObjMap
		}
	}
	return nil
}

// findTargetObjectByOldVal finds the target object by matching oldVal in tfstateSlice
func findTargetObjectByOldVal(oldVal string, tfstateSlice, rawConfigSlice []map[string]interface{}) map[string]interface{} {
	if oldVal == "" {
		return nil
	}

	for _, tfstateObjMap := range tfstateSlice {
		if !ObjectSliceContains(rawConfigSlice, tfstateObjMap) {
			for _, fieldVal := range tfstateObjMap {
				if fmt.Sprintf("%v", fieldVal) == oldVal {
					log.Printf("[DEBUG][findTargetObjectByOldVal] Found object matching oldVal='%s' in tfstateSlice: %v",
						oldVal, tfstateObjMap)
					return tfstateObjMap
				}
			}
		}
	}
	return nil
}

// findTargetObjectFromLocalRemovals finds the target object by matching objectHash fields with localRemovals
func findTargetObjectFromLocalRemovals(d *schema.ResourceData, baseField, objectHash string,
	localRemovals []map[string]interface{}) map[string]interface{} {
	if len(localRemovals) == 0 || d.State() == nil || d.State().Attributes == nil {
		return nil
	}

	// Get all fields for this objectHash from state attributes
	objectFields := make(map[string]interface{})
	for key, val := range d.State().Attributes {
		if strings.HasPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash)) {
			fieldName := strings.TrimPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash))
			objectFields[fieldName] = val
		}
	}

	if len(objectFields) == 0 {
		return nil
	}

	// Try to match them with objects in localRemovals
	for _, localRemovalObj := range localRemovals {
		matches := true
		for fieldName, fieldVal := range objectFields {
			if localRemovalVal, ok := localRemovalObj[fieldName]; !ok || !reflect.DeepEqual(fieldVal, localRemovalVal) {
				matches = false
				break
			}
		}
		if matches {
			log.Printf("[DEBUG][findTargetObjectFromLocalRemovals] Found object by matching objectHash '%s' "+
				"fields with localRemovals: %v", objectHash, localRemovalObj)
			return localRemovalObj
		}
	}
	return nil
}

// findTargetObjectBySingleMatch finds the target object when there's exactly one removed object and one local removal
func findTargetObjectBySingleMatch(tfstateSlice, rawConfigSlice, localRemovals []map[string]interface{}) map[string]interface{} {
	removedObjects := make([]map[string]interface{}, 0)
	for _, tfstateObjMap := range tfstateSlice {
		if !ObjectSliceContains(rawConfigSlice, tfstateObjMap) {
			removedObjects = append(removedObjects, tfstateObjMap)
		}
	}

	// If there's exactly one removed object and one local removal, they should match
	if len(removedObjects) == 1 && len(localRemovals) == 1 {
		if ObjectSliceContains(localRemovals, removedObjects[0]) {
			log.Printf("[DEBUG][findTargetObjectBySingleMatch] Found object by matching single removed object with localRemovals: %v",
				removedObjects[0])
			return removedObjects[0]
		}
	}
	return nil
}

// diffObjectSliceElement handles individual object slice elements
func diffObjectSliceElement(paramKey, oldVal, newVal string, d *schema.ResourceData) bool {
	parts := strings.Split(paramKey, ".")
	if len(parts) < 2 {
		log.Printf("[DEBUG][diffObjectSliceElement] Invalid paramKey format: %s", paramKey)
		return false
	}
	baseField := parts[0]

	// If this is a CREATE scenario (oldVal is empty and newVal has value), don't suppress diff
	if oldVal == "" && newVal != "" {
		log.Printf("[DEBUG][diffObjectSliceElement] CREATE scenario detected (oldVal='', newVal='%s'), "+
			"NOT suppressing diff to show field", newVal)
		return false
	}

	originParamKey := fmt.Sprintf("%s_origin", baseField)
	originVal := d.Get(originParamKey)
	objectHash := parts[1]

	log.Printf("[DEBUG][diffObjectSliceElement] baseField='%s', oldVal='%s', newVal='%s', originVal=%v",
		baseField, oldVal, newVal, originVal)

	// Check if this object is being removed (in localRemovals)
	oldParamVal, newParamVal := d.GetChange(baseField)
	localRemovals := calculateLocalRemovals(originVal, oldParamVal, newParamVal)
	tfstateSlice := convertToObjectSlice(oldParamVal)
	rawConfigSlice := convertToObjectSlice(newParamVal)

	log.Printf("[DEBUG][diffObjectSliceElement] Checking if object '%s' (paramKey='%s') is in localRemovals, localRemovals=%v",
		objectHash, paramKey, formatObjectSliceForLog(localRemovals))

	// Try to find the target object using multiple strategies
	targetObject := findTargetObjectFromState(d, baseField, objectHash, tfstateSlice)
	if targetObject == nil {
		targetObject = findTargetObjectByOldVal(oldVal, tfstateSlice, rawConfigSlice)
	}
	if targetObject == nil {
		targetObject = findTargetObjectFromLocalRemovals(d, baseField, objectHash, localRemovals)
	}
	if targetObject == nil {
		targetObject = findTargetObjectBySingleMatch(tfstateSlice, rawConfigSlice, localRemovals)
	}

	// Check if the target object is in localRemovals
	if targetObject != nil {
		if ObjectSliceContains(localRemovals, targetObject) {
			log.Printf("[DEBUG][diffObjectSliceElement] Object '%s' (paramKey='%s') is in (tfstate - RawConfig) ∩ (origin - RawConfig), "+
				"NOT suppressing diff (local removal)", objectHash, paramKey)
			return false
		}
		log.Printf("[DEBUG][diffObjectSliceElement] Object '%s' (paramKey='%s') is NOT in localRemovals", objectHash, paramKey)
	} else {
		log.Printf("[DEBUG][diffObjectSliceElement] Could not find target object for objectHash '%s' (oldVal='%s'), continuing with normal logic",
			objectHash, oldVal)
	}

	// Handle element removal case
	if newVal == "" {
		log.Printf("[DEBUG][diffObjectSliceElement] Calling handleObjectElementRemoval for paramKey='%s', objectHash='%s'",
			paramKey, objectHash)
		return handleObjectElementRemoval(baseField, objectHash, oldVal, originVal, d)
	}

	// Handle element addition/modification case
	log.Printf("[DEBUG][diffObjectSliceElement] Calling handleObjectElementAddition for paramKey='%s', objectHash='%s'",
		paramKey, objectHash)
	return handleObjectElementAddition(baseField, objectHash, originVal, d)
}

// calculateLocalRemovals calculates local removals: (tfstate - RawConfig) ∩ (origin - RawConfig)
func calculateLocalRemovals(originVal interface{}, oldParamVal, newParamVal interface{}) []map[string]interface{} {
	originSlice := convertToObjectSlice(originVal)
	tfstateSlice := convertToObjectSlice(oldParamVal)
	rawConfigSlice := convertToObjectSlice(newParamVal)

	tfstateMinusRawConfig := FindObjectSliceElementsNotInAnother(tfstateSlice, rawConfigSlice)
	originMinusRawConfig := FindObjectSliceElementsNotInAnother(originSlice, rawConfigSlice)

	localRemovals := make([]map[string]interface{}, 0)
	for _, elem := range tfstateMinusRawConfig {
		if ObjectSliceContains(originMinusRawConfig, elem) {
			localRemovals = append(localRemovals, elem)
		}
	}
	return localRemovals
}

// calculateLocalAdditions calculates local additions: (RawConfig - tfstate) ∩ (RawConfig - origin)
func calculateLocalAdditions(originVal interface{}, oldParamVal, newParamVal interface{}) []map[string]interface{} {
	originSlice := convertToObjectSlice(originVal)
	tfstateSlice := convertToObjectSlice(oldParamVal)
	rawConfigSlice := convertToObjectSlice(newParamVal)

	rawConfigMinusTfstate := FindObjectSliceElementsNotInAnother(rawConfigSlice, tfstateSlice)
	rawConfigMinusOrigin := FindObjectSliceElementsNotInAnother(rawConfigSlice, originSlice)

	localAdditions := make([]map[string]interface{}, 0)
	for _, elem := range rawConfigMinusTfstate {
		if ObjectSliceContains(rawConfigMinusOrigin, elem) {
			localAdditions = append(localAdditions, elem)
		}
	}
	return localAdditions
}

// findObjectByOldVal finds an object in oldObjectList by matching oldVal
func findObjectByOldVal(oldObjectList []interface{}, rawConfigSlice []map[string]interface{}, oldVal string) map[string]interface{} {
	for _, oldItem := range oldObjectList {
		if itemMap, ok := oldItem.(map[string]interface{}); ok {
			if !ObjectSliceContains(rawConfigSlice, itemMap) {
				for _, fieldVal := range itemMap {
					if fmt.Sprintf("%v", fieldVal) == oldVal {
						return itemMap
					}
				}
			}
		}
	}
	return nil
}

// handleObjectElementRemovalWhenOldObjectNil handles the case when oldObject is nil
func handleObjectElementRemovalWhenOldObjectNil(baseField, objectHash, oldVal string, originVal interface{}, d *schema.ResourceData) bool {
	oldParamVal, newParamVal := d.GetChange(baseField)
	if oldParamVal == nil || newParamVal == nil {
		log.Printf("[DEBUG][handleObjectElementRemovalWhenOldObjectNil] Cannot find object '%s' in old state, suppressing diff", objectHash)
		return true
	}

	var oldCount, newCount int
	var oldObjectList []interface{}
	switch v := oldParamVal.(type) {
	case []interface{}:
		oldCount = len(v)
		oldObjectList = v
	case *schema.Set:
		oldCount = v.Len()
		oldObjectList = v.List()
	}
	switch v := newParamVal.(type) {
	case []interface{}:
		newCount = len(v)
	case *schema.Set:
		newCount = v.Len()
	}

	if oldCount <= newCount {
		log.Printf("[DEBUG][handleObjectElementRemovalWhenOldObjectNil] Cannot find object '%s' in old state, suppressing diff", objectHash)
		return true
	}

	var newObjectList []interface{}
	switch v := newParamVal.(type) {
	case []interface{}:
		newObjectList = v
	case *schema.Set:
		newObjectList = v.List()
	}

	localRemovals := calculateLocalRemovals(originVal, oldObjectList, newObjectList)
	log.Printf("[DEBUG][handleObjectElementRemovalWhenOldObjectNil] localRemovals: %v", formatObjectSliceForLog(localRemovals))

	// Try to get the current object by objectHash from state attributes
	currentObject := getObjectFromOldState(d, baseField, objectHash)
	if currentObject == nil {
		currentObject = make(map[string]interface{})
		if d.State() != nil && d.State().Attributes != nil {
			for key, val := range d.State().Attributes {
				if strings.HasPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash)) {
					fieldName := strings.TrimPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash))
					currentObject[fieldName] = val
				}
			}
		}
	}

	// If we still can't get the object, try to find it from oldObjectList using oldVal
	if len(currentObject) == 0 && oldVal != "" {
		rawConfigSlice := convertToObjectSlice(newObjectList)
		currentObject = findObjectByOldVal(oldObjectList, rawConfigSlice, oldVal)
		if currentObject != nil {
			log.Printf("[DEBUG][handleObjectElementRemovalWhenOldObjectNil] Found object '%s' in oldObjectList by matching oldVal='%s': %v",
				objectHash, oldVal, currentObject)
		}
	}

	// Check if the current object is in localRemovals
	if len(currentObject) > 0 {
		if ObjectSliceContains(localRemovals, currentObject) {
			log.Printf("[DEBUG][handleObjectElementRemovalWhenOldObjectNil] Object '%s' (oldVal='%s') is in "+
				"(tfstate - RawConfig) ∩ (origin - RawConfig), NOT suppressing diff (local removal)", objectHash, oldVal)
			return false
		}
		log.Printf("[DEBUG][handleObjectElementRemovalWhenOldObjectNil] Object '%s' (oldVal='%s') is NOT in "+
			"(tfstate - RawConfig) ∩ (origin - RawConfig), suppressing diff (remote addition)", objectHash, oldVal)
		return true
	}

	log.Printf("[DEBUG][handleObjectElementRemovalWhenOldObjectNil] Cannot identify object '%s' (oldVal='%s') by objectHash or oldVal, "+
		"suppressing diff (remote addition)", objectHash, oldVal)
	return true
}

// getFullObjectFromStateAttributes gets the full object from state attributes by objectHash
func getFullObjectFromStateAttributes(d *schema.ResourceData, baseField, objectHash string) map[string]interface{} {
	if d.State() == nil || d.State().Attributes == nil {
		return nil
	}

	fullObject := make(map[string]interface{})
	for key, val := range d.State().Attributes {
		if strings.HasPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash)) {
			fieldName := strings.TrimPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash))
			fullObject[fieldName] = val
		}
	}

	if len(fullObject) == 0 {
		return nil
	}

	log.Printf("[DEBUG][getFullObjectFromStateAttributes] Got full object from state attributes by objectHash '%s': %v",
		objectHash, fullObject)
	return fullObject
}

// findFullObjectByOldVal finds the full object in tfstateSlice by matching oldVal
func findFullObjectByOldVal(oldVal string, tfstateSlice, rawConfigSlice []map[string]interface{}) map[string]interface{} {
	if oldVal == "" {
		return nil
	}

	for _, tfstateObjMap := range tfstateSlice {
		if ObjectSliceContains(rawConfigSlice, tfstateObjMap) {
			continue
		}

		for _, fieldVal := range tfstateObjMap {
			if fmt.Sprintf("%v", fieldVal) == oldVal {
				log.Printf("[DEBUG][findFullObjectByOldVal] Found full object matching oldVal='%s' in tfstateSlice (being removed): %v",
					oldVal, tfstateObjMap)
				return tfstateObjMap
			}
		}
	}

	return nil
}

// findFullObjectByOldObject finds the full object in tfstateSlice by matching oldObject
func findFullObjectByOldObject(oldObject map[string]interface{}, tfstateSlice, rawConfigSlice []map[string]interface{}) map[string]interface{} {
	if len(oldObject) == 0 {
		return nil
	}

	for _, tfstateObjMap := range tfstateSlice {
		if ObjectSliceContains(rawConfigSlice, tfstateObjMap) {
			continue
		}

		matches := true
		for key, val := range oldObject {
			if tfstateVal, ok := tfstateObjMap[key]; !ok || !reflect.DeepEqual(val, tfstateVal) {
				matches = false
				break
			}
		}

		if matches {
			log.Printf("[DEBUG][findFullObjectByOldObject] Found full object matching oldObject in tfstateSlice (being removed): %v",
				tfstateObjMap)
			return tfstateObjMap
		}
	}

	return nil
}

// removalObjectSearchParams contains parameters for searching an object during removal checking
type removalObjectSearchParams struct {
	d              *schema.ResourceData
	baseField      string
	objectHash     string
	oldVal         string
	oldObject      map[string]interface{}
	tfstateSlice   []map[string]interface{}
	rawConfigSlice []map[string]interface{}
}

// findFullObjectForRemoval finds the full object for removal checking
func findFullObjectForRemoval(params removalObjectSearchParams) map[string]interface{} {
	// Strategy 1: Try to get the full object from state attributes by objectHash
	fullObject := getFullObjectFromStateAttributes(params.d, params.baseField, params.objectHash)
	if len(fullObject) > 0 {
		return fullObject
	}

	// Strategy 2: Try to find the full object in tfstateSlice by matching oldVal
	fullObject = findFullObjectByOldVal(params.oldVal, params.tfstateSlice, params.rawConfigSlice)
	if len(fullObject) > 0 {
		return fullObject
	}

	// Strategy 3: Try to find the full object by matching oldObject
	fullObject = findFullObjectByOldObject(params.oldObject, params.tfstateSlice, params.rawConfigSlice)
	return fullObject
}

// checkObjectInLocalRemovals checks if the object is in localRemovals, prioritizing fullObject over oldObject
func checkObjectInLocalRemovals(fullObject, oldObject map[string]interface{},
	localRemovals []map[string]interface{}) bool {
	// Priority 1: Check fullObject if available
	if len(fullObject) > 0 {
		if ObjectSliceContains(localRemovals, fullObject) {
			log.Printf("[DEBUG][checkObjectInLocalRemovals] Full object is in localRemovals")
			return true
		}
		log.Printf("[DEBUG][checkObjectInLocalRemovals] Full object is NOT in localRemovals")
		return false
	}

	// Priority 2: Check oldObject (partial) if fullObject is not available
	if ObjectSliceContains(localRemovals, oldObject) {
		log.Printf("[DEBUG][checkObjectInLocalRemovals] oldObject (partial) is in localRemovals")
		return true
	}

	log.Printf("[DEBUG][checkObjectInLocalRemovals] Could not find full object and oldObject (partial) is NOT in localRemovals")
	return false
}

// handleObjectElementRemoval handles the case when an object element is being removed
func handleObjectElementRemoval(baseField, objectHash, oldVal string, originVal interface{}, d *schema.ResourceData) bool {
	log.Printf("[DEBUG][handleObjectElementRemoval] Object element '%s' is being removed (oldVal='%s'), checking if should suppress diff",
		objectHash, oldVal)

	// Get the object from old state (console value)
	oldObject := getObjectFromOldState(d, baseField, objectHash)

	// If we can't get the object from old state, try to check if it's a remote addition
	if oldObject == nil {
		return handleObjectElementRemovalWhenOldObjectNil(baseField, objectHash, oldVal, originVal, d)
	}

	// Check removal logic: (tfstate - RawConfig) ∩ (origin - RawConfig)
	oldParamVal, newParamVal := d.GetChange(baseField)
	localRemovals := calculateLocalRemovals(originVal, oldParamVal, newParamVal)
	log.Printf("[DEBUG][handleObjectElementRemoval] oldObject=%v, localRemovals=%v", oldObject, formatObjectSliceForLog(localRemovals))

	tfstateSlice := convertToObjectSlice(oldParamVal)
	rawConfigSlice := convertToObjectSlice(newParamVal)

	// Find the full object for checking
	searchParams := removalObjectSearchParams{
		d:              d,
		baseField:      baseField,
		objectHash:     objectHash,
		oldVal:         oldVal,
		oldObject:      oldObject,
		tfstateSlice:   tfstateSlice,
		rawConfigSlice: rawConfigSlice,
	}
	fullObject := findFullObjectForRemoval(searchParams)

	// Check if the object is in localRemovals
	objectInLocalRemovals := checkObjectInLocalRemovals(fullObject, oldObject, localRemovals)

	if objectInLocalRemovals {
		log.Printf("[DEBUG][handleObjectElementRemoval] Object '%s' (oldVal='%s') is in (tfstate - RawConfig) ∩ (origin - RawConfig), "+
			"NOT suppressing diff (local removal)", objectHash, oldVal)
		return false
	}

	log.Printf("[DEBUG][handleObjectElementRemoval] Object '%s' (oldVal='%s') is not in (tfstate - RawConfig) ∩ (origin - RawConfig), "+
		"suppressing diff (remote addition)", objectHash, oldVal)
	return true
}

// findTargetObjectForAddition finds the target object for addition checking
func findTargetObjectForAddition(d *schema.ResourceData, baseField, objectHash string, oldObject map[string]interface{},
	tfstateSlice, localRemovals []map[string]interface{}) map[string]interface{} {
	// First, try to get the object from state attributes by objectHash
	targetObject := findTargetObjectFromState(d, baseField, objectHash, tfstateSlice)
	if targetObject != nil {
		return targetObject
	}

	// If we couldn't find the object from state attributes, try to use oldObject to find it in tfstateSlice
	if len(oldObject) > 0 {
		for _, tfstateObjMap := range tfstateSlice {
			matches := true
			for key, val := range oldObject {
				if tfstateVal, ok := tfstateObjMap[key]; !ok || !reflect.DeepEqual(val, tfstateVal) {
					matches = false
					break
				}
			}
			if matches {
				log.Printf("[DEBUG][findTargetObjectForAddition] Found full object matching oldObject: %v", tfstateObjMap)
				return tfstateObjMap
			}
		}

		// Try direct comparison with oldObject (partial match) against localRemovals
		for _, localRemovalObj := range localRemovals {
			matches := true
			for key, val := range oldObject {
				if localRemovalVal, ok := localRemovalObj[key]; !ok || !reflect.DeepEqual(val, localRemovalVal) {
					matches = false
					break
				}
			}
			if matches {
				log.Printf("[DEBUG][findTargetObjectForAddition] Found matching object in localRemovals by partial oldObject: %v",
					localRemovalObj)
				return localRemovalObj
			}
		}
	}

	return nil
}

// handleObjectElementAdditionWhenNewObjectNil handles the case when newObject is nil
func handleObjectElementAdditionWhenNewObjectNil(baseField, objectHash string, oldObject map[string]interface{},
	originVal interface{}, d *schema.ResourceData) bool {
	log.Printf("[DEBUG][handleObjectElementAdditionWhenNewObjectNil] Cannot find object '%s' in new state (object is in rawconfig but not in state)",
		objectHash)

	// If object is in old state but not in new state, it was removed from script
	if oldObject != nil {
		if !isOriginEmpty(originVal) && isObjectInOrigin(oldObject, originVal) {
			log.Printf("[DEBUG][handleObjectElementAdditionWhenNewObjectNil] Object '%s' exists in old state and was in origin, "+
				"not suppressing diff (allow removal)", objectHash)
			return false
		}
		log.Printf("[DEBUG][handleObjectElementAdditionWhenNewObjectNil] Object '%s' exists in old state but not in new state (remote addition), "+
			"suppressing diff", objectHash)
		return true
	}

	// Try to reconstruct the object from state attributes
	reconstructedObject := getObjectFromState(d, baseField, objectHash)
	if len(reconstructedObject) > 0 {
		if !isOriginEmpty(originVal) && isObjectInOrigin(reconstructedObject, originVal) {
			log.Printf("[DEBUG][handleObjectElementAdditionWhenNewObjectNil] Object '%s' is in rawconfig and origin but not in remote "+
				"state (restoring origin config), not suppressing diff", objectHash)
			return false
		}
	} else if !isOriginEmpty(originVal) {
		log.Printf("[DEBUG][handleObjectElementAdditionWhenNewObjectNil] Object '%s' is in rawconfig but not in remote state, "+
			"and origin is not empty, not suppressing diff (might be restoring origin config)", objectHash)
		return false
	}

	log.Printf("[DEBUG][handleObjectElementAdditionWhenNewObjectNil] Cannot find object '%s' in either new or old state, not suppressing diff",
		objectHash)
	return false
}

// handleObjectElementAddition handles the case when an object element is being added or modified
func handleObjectElementAddition(baseField, objectHash string, originVal interface{}, d *schema.ResourceData) bool {
	// First, check if this object is being removed (in localRemovals)
	oldParamVal, newParamVal := d.GetChange(baseField)
	localRemovals := calculateLocalRemovals(originVal, oldParamVal, newParamVal)
	tfstateSlice := convertToObjectSlice(oldParamVal)

	// Try to get the object from old state to check if it's in localRemovals
	oldObject := getObjectFromOldState(d, baseField, objectHash)
	if oldObject == nil {
		oldObject = make(map[string]interface{})
		if d.State() != nil && d.State().Attributes != nil {
			for key, val := range d.State().Attributes {
				if strings.HasPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash)) {
					fieldName := strings.TrimPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash))
					oldObject[fieldName] = val
				}
			}
		}
	}

	log.Printf("[DEBUG][handleObjectElementAddition] Checking if object '%s' is in localRemovals, oldObject=%v, localRemovals=%v",
		objectHash, oldObject, formatObjectSliceForLog(localRemovals))

	// Try to find the target object to check if it's in localRemovals
	targetObject := findTargetObjectForAddition(d, baseField, objectHash, oldObject, tfstateSlice, localRemovals)

	// Check if the target object is in localRemovals
	if targetObject != nil {
		if ObjectSliceContains(localRemovals, targetObject) {
			log.Printf("[DEBUG][handleObjectElementAddition] Object '%s' is in (tfstate - RawConfig) ∩ (origin - RawConfig), "+
				"NOT suppressing diff (local removal)", objectHash)
			return false
		}
		log.Printf("[DEBUG][handleObjectElementAddition] Object '%s' is NOT in localRemovals", objectHash)
	} else {
		log.Printf("[DEBUG][handleObjectElementAddition] Could not find target object for objectHash '%s', continuing with normal logic",
			objectHash)
	}

	// Try to get the object from new state first
	newObject := getObjectFromState(d, baseField, objectHash)

	// If not found in new state, handle the nil case
	if newObject == nil {
		return handleObjectElementAdditionWhenNewObjectNil(baseField, objectHash, oldObject, originVal, d)
	}

	// If origin is nil or empty, check addition logic
	if isOriginEmpty(originVal) {
		oldParamVal, newParamVal := d.GetChange(baseField)
		localAdditions := calculateLocalAdditions(originVal, oldParamVal, newParamVal)

		if ObjectSliceContains(localAdditions, newObject) {
			log.Printf("[DEBUG][handleObjectElementAddition] Object '%s' is in (RawConfig - tfstate) ∩ (RawConfig - origin) and origin is empty, "+
				"NOT suppressing diff (local addition)", objectHash)
			return false
		}

		log.Printf("[DEBUG][handleObjectElementAddition] Object '%s' is not in (RawConfig - tfstate) ∩ (RawConfig - origin) and origin is empty, "+
			"suppressing diff (remote addition)", objectHash)
		return true
	}

	// Check if this object is in origin
	if isObjectInOrigin(newObject, originVal) {
		return handleObjectElementAdditionWhenInOrigin(baseField, objectHash, newObject, oldObject, d)
	}

	// Check addition logic: (RawConfig - tfstate) ∩ (RawConfig - origin)
	localAdditions := calculateLocalAdditions(originVal, oldParamVal, newParamVal)

	if ObjectSliceContains(localAdditions, newObject) {
		log.Printf("[DEBUG][handleObjectElementAddition] Object '%s' is in (RawConfig - tfstate) ∩ (RawConfig - origin), "+
			"NOT suppressing diff (local addition)", objectHash)
		return false
	}

	log.Printf("[DEBUG][handleObjectElementAddition] Object '%s' is not in (RawConfig - tfstate) ∩ (RawConfig - origin), "+
		"suppressing diff (remote addition)", objectHash)
	return true
}

// handleObjectElementAdditionWhenInOrigin handles the case when the object is in origin
func handleObjectElementAdditionWhenInOrigin(baseField, objectHash string, newObject, oldObject map[string]interface{}, d *schema.ResourceData) bool {
	// Check if it exists in remote state by matching object content
	oldParamVal, _ := d.GetChange(baseField)
	var oldObjectList []interface{}
	if oldParamVal != nil {
		switch v := oldParamVal.(type) {
		case []interface{}:
			oldObjectList = v
		case *schema.Set:
			oldObjectList = v.List()
		}
	}

	// Check if newObject exists in remote state by content matching
	objectInRemoteState := false
	for _, oldItem := range oldObjectList {
		if oldItemMap, ok := oldItem.(map[string]interface{}); ok {
			if reflect.DeepEqual(oldItemMap, newObject) {
				objectInRemoteState = true
				log.Printf("[DEBUG][handleObjectElementAdditionWhenInOrigin] Object '%s' found in remote state by content matching: %v",
					objectHash, newObject)
				break
			}
		}
	}

	if !objectInRemoteState {
		log.Printf("[DEBUG][handleObjectElementAdditionWhenInOrigin] Object '%s' (content: %v) is in rawconfig and origin but not in remote "+
			"state (restoring origin config), not suppressing diff",
			objectHash, newObject)
		return false
	}

	// Object exists in both origin and remote state
	if oldObject != nil && reflect.DeepEqual(oldObject, newObject) {
		log.Printf("[DEBUG][handleObjectElementAdditionWhenInOrigin] Object '%s' unchanged and in origin and remote state, "+
			"not suppressing diff to preserve config value", objectHash)
		return false
	}

	log.Printf("[DEBUG][handleObjectElementAdditionWhenInOrigin] Object '%s' was in origin and remote state but changed, "+
		"suppressing diff", objectHash)
	return true
}

// isObjectInOrigin checks if an object exists in origin value
func isObjectInOrigin(obj map[string]interface{}, originVal interface{}) bool {
	if originVal == nil {
		return false
	}

	var originList []interface{}
	switch v := originVal.(type) {
	case []interface{}:
		originList = v
	case *schema.Set:
		originList = v.List()
	default:
		return false
	}

	for _, item := range originList {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if reflect.DeepEqual(itemMap, obj) {
				return true
			}
		}
	}

	return false
}

// getObjectFromState retrieves an object from the resource state by its hash
func getObjectFromState(d *schema.ResourceData, baseField, objectHash string) map[string]interface{} {
	// For TypeSet, we need to reconstruct the object from individual fields
	// Since TypeSet uses hash, we can't directly access by hash
	// Instead, we need to reconstruct the object from individual fields
	obj := make(map[string]interface{})
	hasFields := false

	// Try to get all fields of this object from the state
	// The paramKey format for object fields is: baseField.objectHash.fieldName
	if d.State() != nil && d.State().Attributes != nil {
		for key := range d.State().Attributes {
			if strings.HasPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash)) {
				fieldName := strings.TrimPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash))
				obj[fieldName] = d.Get(key)
				hasFields = true
			}
		}
	}

	if !hasFields {
		return nil
	}

	return obj
}

// getObjectFromOldState retrieves an object from the old state (console value) by its hash
func getObjectFromOldState(d *schema.ResourceData, baseField, objectHash string) map[string]interface{} {
	// First, try to get the object from state attributes (for old state)
	obj := make(map[string]interface{})
	hasFields := false

	// Try to get all fields of this object from the old state attributes
	// The paramKey format for object fields is: baseField.objectHash.fieldName
	if d.State() != nil && d.State().Attributes != nil {
		for key, val := range d.State().Attributes {
			if strings.HasPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash)) {
				fieldName := strings.TrimPrefix(key, fmt.Sprintf("%s.%s.", baseField, objectHash))
				// In diff suppression, we need to check if this is from old state
				// The old state value might be in the state attributes
				obj[fieldName] = val
				hasFields = true
			}
		}
	}

	// If we found fields in state, return the object
	if hasFields {
		return obj
	}

	// If not found in state attributes, try to reconstruct from GetChange old value
	// Since TypeSet uses hash, we need to find the object by matching all fields
	oldParamVal, _ := d.GetChange(baseField)
	if oldParamVal != nil {
		var oldObjectList []interface{}
		switch v := oldParamVal.(type) {
		case []interface{}:
			oldObjectList = v
		case *schema.Set:
			oldObjectList = v.List()
		}

		// Try to find the object by checking if any object in old value has fields matching the hash
		// We can't directly match by hash, but we can try to find objects that are not in new value
		// For now, return nil and let the caller handle it based on count comparison
		_ = oldObjectList
	}

	// If not found, return nil and let the caller handle it.
	// The caller will check if the object is in origin, and if not, assume it's a remote addition.
	return nil
}

// checkObjectInRemoteState checks if an object exists in remote state (console value)
func checkObjectInRemoteState(baseField string, obj map[string]interface{}, d *schema.ResourceData) bool {
	// Get the old value (console/remote state) from GetChange
	oldParamVal, _ := d.GetChange(baseField)
	if oldParamVal == nil {
		return false
	}

	var objectList []interface{}
	switch v := oldParamVal.(type) {
	case []interface{}:
		objectList = v
	case *schema.Set:
		objectList = v.List()
	default:
		return false
	}

	// Check if the object exists in the list
	for _, item := range objectList {
		if itemMap, ok := item.(map[string]interface{}); ok {
			if reflect.DeepEqual(itemMap, obj) {
				log.Printf("[DEBUG][checkObjectInRemoteState] Object already exists in remote state, suppressing diff")
				return true
			}
		}
	}

	log.Printf("[DEBUG][checkObjectInRemoteState] Object does not exist in remote state, allowing change")
	return false
}

// convertToObjectSlice converts interface{} to []map[string]interface{}
func convertToObjectSlice(val interface{}) []map[string]interface{} {
	if val == nil {
		return nil
	}

	var result []map[string]interface{}
	switch v := val.(type) {
	case []interface{}:
		for _, item := range v {
			if itemMap, ok := item.(map[string]interface{}); ok {
				result = append(result, itemMap)
			}
		}
	case *schema.Set:
		for _, item := range v.List() {
			if itemMap, ok := item.(map[string]interface{}); ok {
				result = append(result, itemMap)
			}
		}
	}

	return result
}

// diffObjectSliceParam is used to check whether the parameters of the current object slice type have been modified
// other than those changed in the console.
// Only show diff in the following two scenarios (return false to not suppress diff):
//  1. New additions: (RawConfig - tfstate) ∩ (RawConfig - origin)
//     Elements that are in RawConfig but not in tfstate AND not in origin
//  2. Removals: (tfstate - RawConfig) ∩ (origin - RawConfig)
//     Elements that are in tfstate but not in RawConfig AND in origin but not in RawConfig
//
// All other scenarios will suppress the diff (return true).
func diffObjectSliceParam(paramKey string, d *schema.ResourceData) bool {
	originParamKey := fmt.Sprintf("%s_origin", paramKey)
	originVal := d.Get(originParamKey)
	originSlice := convertToObjectSlice(originVal)

	// Get old and new values from GetChange
	oldParamVal, newParamVal := d.GetChange(paramKey)
	tfstateSlice := convertToObjectSlice(oldParamVal)   // tfstate (remote state)
	rawConfigSlice := convertToObjectSlice(newParamVal) // RawConfig (script config)

	log.Printf("[DEBUG][diffObjectSliceParam] paramKey='%s', originSlice length=%d, tfstateSlice length=%d, rawConfigSlice length=%d",
		paramKey, len(originSlice), len(tfstateSlice), len(rawConfigSlice))

	// Scenario 1: New additions - (RawConfig - tfstate) ∩ (RawConfig - origin)
	localAdditions := calculateLocalAdditions(originVal, oldParamVal, newParamVal)
	if len(localAdditions) > 0 {
		log.Printf("[DEBUG][diffObjectSliceParam] Found local additions (RawConfig - tfstate) ∩ (RawConfig - origin): %d, not suppressing diff",
			len(localAdditions))
		return false
	}

	// Scenario 2: Removals - (tfstate - RawConfig) ∩ (origin - RawConfig)
	localRemovals := calculateLocalRemovals(originVal, oldParamVal, newParamVal)
	if len(localRemovals) > 0 {
		log.Printf("[DEBUG][diffObjectSliceParam] Found local removals (tfstate - RawConfig) ∩ (origin - RawConfig): %d, not suppressing diff",
			len(localRemovals))
		return false
	}

	// No local additions or removals, suppress the diff
	log.Printf("[DEBUG][diffObjectSliceParam] No local additions or removals, suppressing diff")
	return true
}

// FindObjectSliceElementsNotInAnother returns elements from source that are not in target
// This is equivalent to source - target (set difference) for object slices
func FindObjectSliceElementsNotInAnother(source, target []map[string]interface{}) []map[string]interface{} {
	var result []map[string]interface{}
	for _, sv := range source {
		if !ObjectSliceContains(target, sv) {
			result = append(result, sv)
		}
	}
	return result
}

// ObjectSliceContains checks if a target object is present in a slice of objects
func ObjectSliceContains(slice []map[string]interface{}, target map[string]interface{}) bool {
	for _, v := range slice {
		if reflect.DeepEqual(v, target) {
			return true
		}
	}
	return false
}

// formatObjectSliceForLog formats an object slice for logging purposes
func formatObjectSliceForLog(slice []map[string]interface{}) string {
	if len(slice) == 0 {
		return "[]"
	}
	var result []string
	for _, obj := range slice {
		// Try to extract name field if it exists, otherwise use the whole object
		if name, ok := obj["name"].(string); ok {
			result = append(result, name)
		} else {
			result = append(result, fmt.Sprintf("%v", obj))
		}
	}
	return fmt.Sprintf("[%s]", strings.Join(result, ", "))
}
