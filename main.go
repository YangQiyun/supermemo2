package main

import (
	"math"
	"errors"
	"strconv"
	"fmt"
)
//ref
//https://github.com/sunaiwen/supermemo2.js/blob/master/test/dist.js
//https://github.com/tqyq/supermemo2-java-demo/blob/master/sm2/Test.java#L42
//sm https://www.supermemo.com/english/ol/sm2.htm
//sm2_plus http://www.blueraja.com/blog/477/a-better-spaced-repetition-learning-algorithm-sm2

/*
quality:
5 - perfect response
4 - correct response after a hesitation
3 - correct response recalled with serious difficulty
2 - incorrect response; where the correct one seemed easy to recall
1 - incorrect response; the correct one remembered
0 - complete blackout.

*/

type Item struct {
	ef      float64
	n       int //interval day
	quality int //every review quality from [0-5]
}

const (
	DEFAULT_EF        = 2.5
	DEFAULT_N         = 1
	DEFAULT_Q         = 3
	INVALID_PARAMETER = "invalid quality"
)

func GetNewItem() (*Item) {
	item := &Item{
		ef:      DEFAULT_EF,
		n:       DEFAULT_N,
		quality: 3,
	}

	return item
}

func (item *Item) SM2(q int) (error) {

	if q < 0 || q > 5 {
		return errors.New(INVALID_PARAMETER + " " + strconv.Itoa(q))
	}
	floatQ := float64(q)
	//If the quality response was lower than 3 then start repetitions for the item from the beginning without changing the E-Factor
	if q < DEFAULT_Q {
		item.ef = DEFAULT_EF
		item.n = DEFAULT_N
		return nil
	}
	//EF':=EF+(0.1-(5-q)*(0.08+(5-q)*0.02))
	item.ef = item.ef + (0.1 - (5-floatQ)*(0.08+(5-floatQ)*0.02))
	//If EF is less than 1.3 then let EF be 1.3.
	if item.ef < 1.3 {
		item.ef = 1.3
	}
	//I(1):=1
	//I(2):=6
	//for n>2: I(n):=I(n-1)*EF
	switch item.n {
	case 1:
		item.n = 6
		item.ef = DEFAULT_EF
	default:
		item.n = int(math.Floor(float64(item.n) * item.ef))
	}
	return nil
}

//difficulty [0.0,1.0]
//performance [0.0,1.0]

const (
	DEFAULT_DIFFICULTY  = 0.3
	DEFAULT_PERFORMANCE = 0.6
)

type Item_Plus struct {
	dateLastReviewed   int
	percentOverdue     float64
	difficulty         float64
	difficultyWeight   float64
	daysBetweenReviews float64
}

func GetNewItemPlus(newDate int) (*Item_Plus) {
	item := &Item_Plus{
		dateLastReviewed:   1,
		percentOverdue:     1,
		difficulty:         DEFAULT_DIFFICULTY,
		difficultyWeight:   3 - 1.7*DEFAULT_DIFFICULTY,
		daysBetweenReviews: 1,
	}
	return item
}

func (item *Item_Plus) SM2_PLUS(nowDate int,performance float64) (error) {
	var (
		err error
	)
	if performance < 0 || performance > 1 {
		return errors.New(INVALID_PARAMETER + " " + strconv.FormatFloat(performance, 'f', 6, 64))
	}

	if performance < DEFAULT_PERFORMANCE {
		item.percentOverdue = 1
		item.dateLastReviewed = nowDate
		if 1/(item.difficultyWeight*item.difficultyWeight) < 1 {
			item.daysBetweenReviews = 1
		}
	}

	item.percentOverdue = math.Min(2, float64((nowDate - item.dateLastReviewed))/item.daysBetweenReviews)
	item.difficulty += item.percentOverdue * (8 - 9*performance) / 17
	if item.difficulty < 0 {
		item.difficulty = 0
	}
	if item.difficulty > 1 {
		item.difficulty = 1
	}

	item.difficultyWeight = 3 - 1.7*item.difficulty
	item.daysBetweenReviews *= 1 + (item.difficultyWeight-1)*item.percentOverdue
	item.dateLastReviewed = nowDate
	return err
}

func TestSM2() {

	item := GetNewItem()
	fmt.Printf("init %v\n",*item)
	fmt.Println("q=2")
	item.SM2(2)
	fmt.Printf("%v\n",*item)
	fmt.Println("q=4")
	item.SM2(4)
	fmt.Printf("%v\n",*item)
	fmt.Println("q=5")
	item.SM2(5)
	fmt.Printf("%v\n",*item)
	fmt.Println("q=3")
	item.SM2(3)
	fmt.Printf("%v\n",*item)
}


func TestSM2PLUS() {
	nowdate := 1
	item := GetNewItemPlus(nowdate)
	fmt.Printf("init %v\n",*item)
	fmt.Println("q=0.7")
	nowdate += int(math.Round(item.daysBetweenReviews))
	item.SM2_PLUS(nowdate,0.7)
	fmt.Printf("day is %d   %v\n",nowdate,*item)
	fmt.Println("q=0.2")
	nowdate += int(math.Round(item.daysBetweenReviews))
	item.SM2_PLUS(nowdate,0.2)
	fmt.Printf("day is %d   %v\n",nowdate,*item)
	fmt.Println("q=1")
	nowdate += int(math.Round(item.daysBetweenReviews))
	item.SM2_PLUS(nowdate,1)
	fmt.Printf("day is %d   %v\n",nowdate,*item)
	fmt.Println("q=1")
	nowdate += int(math.Round(item.daysBetweenReviews))
	item.SM2_PLUS(nowdate,1)
	fmt.Printf("day is %d   %v\n",nowdate,*item)
}

func SM2PLUS_ACTIVE()  {
	nowdate := 1
	item := GetNewItemPlus(nowdate)
	item.SM2_PLUS(2,1)
	fmt.Println(item.daysBetweenReviews)
	item.SM2_PLUS(3,1)
	fmt.Println(item.daysBetweenReviews)
	item.SM2_PLUS(4,1)
	fmt.Println(item.daysBetweenReviews)
	item.SM2_PLUS(5,1)
	fmt.Println(item.daysBetweenReviews)

}

func main()  {
	TestSM2()
	fmt.Println("-----------------")
	TestSM2PLUS()
	fmt.Println("-----------------")
	SM2PLUS_ACTIVE()

}
