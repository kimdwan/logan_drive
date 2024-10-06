import { useEffect, useState } from "react"

// 유저의 컴퓨터 넘버를 불러오는 함수
export const LoadComputerNumber = () => {
  const [ computerNumber, setComputerNumber ] = useState("")

  useEffect(() => {
    const localstrage_computer_number = localStorage.getItem("logan_computer_number")
    if (localstrage_computer_number !== "") {
      setComputerNumber(localstrage_computer_number)
    }

  }, [])

  return { computerNumber, setComputerNumber }
}