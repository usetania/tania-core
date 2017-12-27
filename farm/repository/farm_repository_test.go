package repository

import (
	"fmt"
	"testing"

	"github.com/Tanibox/tania-server/farm/entity"
	"github.com/stretchr/testify/assert"
)

func TestFarmInMemorySave(t *testing.T) {
	// Given
	done := make(chan bool)
	repo := NewFarmRepositoryInMemory()

	farm1, _ := entity.CreateFarm("My Farm Family", "", "-90.000", "-180.000", "organic", "ID", "JK")
	farm2, _ := entity.CreateFarm("My Second Farm", "", "-90.000", "-180.000", "organic", "ID", "JK")

	fmt.Println(farm2)

	// When
	var saveResult1, saveResult2, count1 RepositoryResult
	go func() {
		saveResult1 = <-repo.Save(&farm1)
		saveResult2 = <-repo.Save(&farm2)

		count1 = <-repo.Count()
		done <- true
	}()

	// Then
	<-done
	assert.NotNil(t, saveResult1)

	assert.Equal(t, count1.Result, 2)

}
