package assert

import "testing"

func Equal(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) bool {
    if expected != actual {
        t.Errorf("not equal: expected %v got %v", expected, actual)
        return false
    }
    return true
}
