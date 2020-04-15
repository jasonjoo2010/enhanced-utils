package strutils

import (
	"testing"

	"gotest.tools/assert"
)

func TestToUnderscore(t *testing.T) {
	assert.Equal(t, "", ToUnderscore(""))
	assert.Equal(t, "", ToUnderscore(" "))
	assert.Equal(t, "_", ToUnderscore(" _"))
	assert.Equal(t, "a_apple_banana_b", ToUnderscore(" aAppleBanana_b "))
	assert.Equal(t, "a_apple_banana_b_", ToUnderscore("aAppleBanana_b_"))
	assert.Equal(t, "_a_apple_banana_b_", ToUnderscore("_aAppleBanana_b_"))
	assert.Equal(t, "apple_banana", ToUnderscore("AppleBanana"))
	assert.Equal(t, "___a_apple_____banana____b____", ToUnderscore("__AApple____Banana____b____"))
}

func TestToCamel(t *testing.T) {
	assert.Equal(t, "", ToCamel(""))
	assert.Equal(t, "", ToCamel(" "))
	assert.Equal(t, "", ToCamel(" _"))
	assert.Equal(t, "aBCD", ToCamel(" _a_b_c_d"))
	assert.Equal(t, "appleBanana", ToCamel("AppleBanana"))
	assert.Equal(t, "appleBanana", ToCamel("apple_banana"))
	assert.Equal(t, "appleBanana", ToCamel("_apple_banana"))
	assert.Equal(t, "appleBanana", ToCamel("apple_banana"))
	assert.Equal(t, "appleBanana", ToCamel("_Apple_Banana"))
}
