package shorter

import (
	"errors"
	"math"
)

var decimalMap map[uint64]rune
var revDecimalMap map[rune]uint64
func init() {
	decimalMap = make(map[uint64]rune)
	revDecimalMap = make(map[rune]uint64)
	var zero = '0'
	var runesOfDigitsLowersUppers []rune
	for i:=0;i<=9;i++ {
		runesOfDigitsLowersUppers = append(runesOfDigitsLowersUppers,zero+rune(i))
	}
	var lowerA = 'a'
	for i:=0;i<26;i++ {
		runesOfDigitsLowersUppers = append(runesOfDigitsLowersUppers,lowerA+rune(i))
	}
	var upperA = 'A'
	for i := 0; i < 26; i++ {
		runesOfDigitsLowersUppers = append(runesOfDigitsLowersUppers,upperA+rune(i))

	}
	for index,value := range runesOfDigitsLowersUppers {
		decimalMap[uint64(index)] = value
		revDecimalMap[value] = uint64(index)
	}
}

func convert10To62(num uint64)string {
	if num == 0 {
		return "0"
	}
	var res []rune
	for num != 0 {
		res =append(res, decimalMap[num-(num/62)*62])
		num = num/62
	}
	for i:=0;i<len(res)/2;i++ {
		res[i],res[len(res)-i-1] =res[len(res)-i-1],res[i]
	}
	return string(res)

}
const mayOverFlow = (math.MaxUint64/62)+1
func convert62To10(str string) (uint64,error) {
	slice := []rune(str)
	var needReturn uint64

	for j:=0;j<len(slice);j++ {
		index,ok := revDecimalMap[slice[j]]
		if !ok {
			return 0,errors.New("can't convert invalid data with "+str)
		}
		if needReturn >= mayOverFlow {
			return 0, errors.New("may overflow")
		}
		needReturn *= 62
		n1 := needReturn+index
		if n1 < needReturn || n1 > math.MaxUint64 {
			return 0, errors.New("may overflow")
		}
		needReturn = n1

	}

	return needReturn,nil
}

