package domain

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateInventory(t *testing.T) {
	// When
	inventoryMaterial1, err1 := CreateInventoryMaterial(Vegetable{}, "Sawi Putih")
	inventoryMaterial2, err2 := CreateInventoryMaterial(Fruit{}, "Apple")

	// Then
	assert.Nil(t, err1)
	assert.NotNil(t, inventoryMaterial1)
	assert.Nil(t, err2)
	assert.NotNil(t, inventoryMaterial2)
}
