package main

import (
	"fmt"
)

// func main() {
// 	//fmt.Println("Hello Go")
// 	var i int
// 	var j float32 = 5635.23

// 	i = 42
// 	fmt.Println(i)
// 	fmt.Println(j)

// 	firstName := "virola"
// 	fmt.Println(firstName)

// 	b := true

// 	fmt.Println(b)

// 	c := complex(2, 3)

// 	fmt.Println(c)

// 	r, im := real(c), imag(c)
// 	println(r, im)
// }

// func main() {

// 	var firstName *string = new(string)
// 	*firstName = "Virolita"
// 	fmt.Println(*firstName)

// 	secondName := "Ale"
// 	fmt.Println(secondName)

// 	ptr := &secondName
// 	fmt.Println(ptr, *ptr)

// 	var intPtr *int = new(int)
// 	*intPtr = 45
// 	fmt.Println(intPtr, *intPtr)

// 	const con = "3.2"
// 	fmt.Println(con)

// }

// func main() {
// 	arr := [3]int{1, 2, 3}
// 	//var arr [3]int
// 	// arr[0] = 1
// 	// arr[1] = 2
// 	// arr[2] = 3

// 	slice := arr[:]

// 	arr[1] = 42
// 	arr[2] = 27
// 	fmt.Println("slice ", slice)
// 	fmt.Println(arr)
// }

// func main() {
// 	slice := []int{1, 2, 3}

// 	slice = append(slice, 4, 5, 6)

// 	fmt.Println(slice)

// 	s2 := slice[1:]
// 	s3 := slice[:1]
// 	s4 := slice[1:2]

// 	fmt.Println("s2", s2, "s3", s3, "s4", s4)
// }

// func main() {

// 	m := map[string]int{"foo": 42}

// 	fmt.Println(m["foo"])

// 	m["foo"] = 46

// 	fmt.Println(m["foo"])

// 	delete(m, "foo")

// 	fmt.Println(m["foo"])
// }

//

//"github.com/pluralsight/webservice/models"

// func main() {
// 	u := &models.User{
// 		Id:        2,
// 		FirstName: "Pepe",
// 		LastName:  "Virola",
// 	}

// 	u2 := &models.User{
// 		Id:        3,
// 		FirstName: "Pepe",
// 		LastName:  "Virola",
// 	}

// 	users := []*models.User{u}
// 	users = append(users, u2)

// 	fmt.Println("Users", *users[0])

// }

// func main() {

// 	port := 3000

// 	_, err := startWebServer(port, 3)
// 	fmt.Println(err)

// }

// func startWebServer(port, numberOfRetries int) (int, error) {

// 	test := false

// 	fmt.Println("Starting server....")
// 	fmt.Println("Port", port, numberOfRetries)
// 	fmt.Println("Server started")

// 	if !test {
// 		return port, nil
// 	} else {
// 		return port, errors.New("something went wrong")
// 	}

// }

// func main() {
// 	controllers.RegisterControllers()
// 	http.ListenAndServe(":3000", nil)
// 	fmt.Println("Hello Go")

// }

// func twoSum(nums []int, target int) ([]int, error) {

// 	if len(nums) < 2 || len(nums) > int(math.Pow(10, 4)) {
// 		return nil, errors.New("sorry invalid number")
// 	}

// 	result := []int{}
// 	var sum int
// 	for i := 0; i < len(nums); i++ {
// 		value := nums[i]
// 		index := i
// 		for j := 0; j < len(nums); j++ {
// 			if index != j {
// 				sum = value + nums[j]
// 				if sum == target {
// 					result = append(result, i)
// 					result = append(result, j)
// 					break
// 				}
// 			}
// 		}
// 		if sum == target {
// 			break
// 		}
// 	}

// 	return result, nil
// }

// func twoSum(nums []int, target int) ([]int, error) {

// 	if len(nums) < 2 || len(nums) > int(math.Pow(10, 4)) {
// 		return nil, errors.New("sorry invalid number")
// 	}

// 	var seen map[int]int
// 	result := []int{}
// 	//var sum int
// 	seen = make(map[int]int)

// 	for i, num := range nums {
// 		complement := target - num
// 		if j, ok := seen[complement]; ok {
// 			result = []int{i, j}
// 			break
// 		}
// 		seen[num] = i
// 	}
// 	//fmt.Println(seen)

// 	return result, nil
// }

// func main() {

// 	num := []int{3, 2, 4}
// 	var target int
// 	target = 6
// 	res, err := fmt.Println(twoSum(num, target))
// 	if err != nil {
// 		fmt.Println(res)
// 	} else {
// 		fmt.Println(err)
// 	}

// }

type ListNode struct {
	Val int
	Next *ListNode
}

/**
 * Definition for singly-linked list.
 * type ListNode struct {
 *     Val int
 *     Next *ListNode
 * }
 */
func addTwoNumbers(l1 *ListNode, l2 *ListNode) *ListNode {

		
	//reverseL1 := reverseList(l1)
	//reverseL2 := reverseList(l2)
	reverseL1 := l1
	reverseL2 := l2
	var sumValueL1 [] int
	var sumValueL2 [] int
	var sumResult int
	var sumResultAsSlice [] int

	for node := reverseL1; node != nil; node = node.Next {
		//fmt.Println(node.val)		
		sumValueL1 = append(sumValueL1, node.Val)
		fmt.Println(sumValueL1)

	}

	for node := reverseL2; node != nil; node = node.Next {
		sumValueL2 = append(sumValueL2, node.Val)
		fmt.Println(sumValueL2)

	}

	sumResult = sliceToNumber(sumValueL1) + sliceToNumber(sumValueL2) 
    sumResultAsSlice = numberToSlice(sumResult)
		

	//return reverseList(createListFromSlice(sumResultAsSlice))
	return createListFromSlice(sumResultAsSlice)
    
}

func sliceToNumber(nums []int) int {
    result := 0
    for _, num := range nums {
        result = result*10 + num
    }
    return result
}

func numberToSlice(num int) [] int {
    if num == 0 {
        return []int{0}
    }
    var result []int
    for num > 0 {
        digit := num % 10
        result = append([]int{digit}, result...) // Prepend digit
        num /= 10
    }
    return result

}


func reverseList(head *ListNode) *ListNode{

	var prev *ListNode	
	var curr *ListNode = head

	for curr != nil {		
		next := curr.Next // Save the next node
		curr.Next = prev  // Reverse the link
		prev = curr		  // Move prev forward
		curr = next		  // Move curr forward
	}
	
	return prev

}


func createListFromSlice(nums []int) *ListNode {
    if len(nums) == 0 {
        return nil
    }
    
    head := &ListNode{Val: nums[0]}
    current := head
    
    for i := 1; i < len(nums); i++ {
        current.Next = &ListNode{Val: nums[i]}
        current = current.Next
    }
    
    return head
}

func main(){
	
 
	 l1 :=createListFromSlice([]int{2,4,3})
	 l2 :=createListFromSlice([]int{5,6,4})

	//l1 :=createListFromSlice([]int{1,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,1})
	//l2 :=createListFromSlice([]int{5,6,4})


	var result = addTwoNumbers(l1, l2)
	for node := result; node != nil; node = node.Next {
		fmt.Print(node.Val, ",")	
	}

	
}

// func main() {

// 	for i := 0; i <= 5; i++ {

// 		if i == 4 {
// 			continue
// 		} else {
// 			fmt.Println(i)
// 		}

// 	}

// }

// func main() {

// 	slice := []int{1, 2, 3}

// 	for _, v := range slice {
// 		fmt.Println(v)

// 	}

// }

// func main() {

// 	fmt.Println("start")
// 	//panic("something went wrong ")
// 	fmt.Println("end")

// }

// type User struct {
// 	Id        int
// 	FirstName string
// 	LastName  string
// }

// func main() {
// 	u1 := User{
// 		Id:        1,
// 		FirstName: "Pepe",
// 		LastName:  "Virola",
// 	}

// 	u2 := User{
// 		Id:        1,
// 		FirstName: "Ozzy",
// 		LastName:  "Obsburne",
// 	}

// 	if u1 == u2 {
// 		println("Same user")
// 	} else if u1.Id == u2.Id {
// 		println("Similar user")
// 		println("u1", &u1)
// 		println("u2", &u2)
// 	} else {
// 		println("Different user")
// 	}

// }
