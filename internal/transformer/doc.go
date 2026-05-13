// Package transformer applies structural transformations to .env key-value maps.
//
// Supported transformations include:
//   - Prefix: prepend a string to all keys
//   - Suffix: append a string to all keys
//   - RenameMap: rename specific keys by mapping old name to new name
//   - UppercaseKeys: convert all keys to uppercase
//   - LowercaseKeys: convert all keys to lowercase
//
// Transformations are applied in the following order:
//  1. Rename
//  2. Case conversion
//  3. Prefix / Suffix
//
// Values are never modified; only keys are affected.
package transformer
