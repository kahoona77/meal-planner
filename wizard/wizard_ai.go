package wizard

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"github.com/sjwhitworth/golearn/base"
	"github.com/sjwhitworth/golearn/evaluation"
	"github.com/sjwhitworth/golearn/trees"
	"io"
	"meal-planner/planner"
)

type AiWizard struct {
}

func (aiw *AiWizard) Generate(wizardWeek Week) (*planner.Week, error) {

	// Convert past weeks of meals to a format compatible with golearn
	data := [][]string{
		{"Date", "MealId", "Name"},
		{"2024-02-01", "30", "Wurst"},
		{"2024-02-02", "5", "Mais"},
		{"2024-02-03", "13", "Burger"},
		{"2024-02-04", "8", "Erbse"},
		{"2024-02-05", "22", "Möhre"},
		{"2024-02-06", "4", "Müsli"},
		{"2024-02-07", "32", "Fleisch"},
		{"2024-02-08", "11", "Nudeln"},
		{"2024-02-09", "7", "Brot"},
		{"2024-02-10", "26", "Toast"},
		{"2024-02-11", "22", "Möhre"},
		{"2024-02-12", "19", "Pizza"},
		{"2024-02-13", "4", "Müsli"},
		{"2024-02-14", "1", "Reis"},
	}
	reader := generateCSVReader(data)

	//// Define attributes
	//attributes := []base.Attribute{
	//	base.NewFloatAttribute("Date"),
	//	base.NewFloatAttribute("MealID"),
	//}

	// Create a new Instances object
	instances, err := base.ParseCSVToInstancesFromReader(reader, true)
	if err != nil {
		return nil, err
	}

	// Create a decision tree classifier
	tree := trees.NewID3DecisionTree(1)

	// Train the classifier
	trainData, testData := base.InstancesTrainTestSplit(instances, 0.60)
	err = tree.Fit(trainData)
	if err != nil {
		return nil, err
	}

	// Predict meals for the upcoming week
	//futureWeek := &planner.Week{
	//	Start:  time.Now().AddDate(0, 0, 7),  // Next week
	//	End:    time.Now().AddDate(0, 0, 13), // End of next week
	//	Offset: 7,
	//	Meals:  make([]*planner.MealOfTheDay, 7), // Placeholder for predicted meals
	//}

	predictions, err := tree.Predict(testData)
	if err != nil {
		return nil, err
	}

	// Evaluate
	fmt.Println("ID3 Performance (information gain)")
	cf, err := evaluation.GetConfusionMatrix(testData, predictions)
	if err != nil {
		panic(fmt.Sprintf("Unable to get confusion matrix: %s", err.Error()))
	}
	fmt.Println(evaluation.GetSummary(cf))

	//for i := 0; i < 7; i++ {
	//	date := futureWeek.Start.AddDate(0, 0, i)
	//
	//	mealID, _ := base.GetClass(prediction)
	//	futureWeek.Meals[i] = &planner.MealOfTheDay{
	//		Id:   fmt.Sprintf("PredictedMeal%d", i+1),
	//		Date: date,
	//		//MealId: int64(mealID),
	//	}
	//}

	return nil, fmt.Errorf("implemnt me")
}

// Create a function to generate an io.Reader backed by CSV data
func generateCSVReader(data [][]string) io.ReadSeeker {
	// Create a buffer to store CSV data
	var buf bytes.Buffer

	// Create a CSV writer
	writer := csv.NewWriter(&buf)

	// Write the data to the CSV writer
	err := writer.WriteAll(data)
	if err != nil {
		panic(err)
	}

	// Return the buffer as an io.Reader
	return bytes.NewReader(buf.Bytes())
}
