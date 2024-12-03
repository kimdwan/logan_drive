

export class ContentsNavbarUserProfileFunction {

  constructor (computerNumber, setComputerNumber, navigate) {
    this.computerNumber = computerNumber
    this.setComputerNumber = setComputerNumber
    this.navigate = navigate
  }

  async GetUserData (url) {

    try {
      const response = await fetch(url, {
        method : "GET",
        headers : {
          "Content-Type" : "application/json; charset=utf-8",
          "X-Requested-With" : "XMLHttpRequest",
          "User-Computer-Number" : this.computerNumber,
        },
        credentials : "include",
      })

      if (!response.ok) {
        if (response.status === 401) {
          alert("세션이 만료되었습니다 다시 로그인 해주세요")
          this.setComputerNumber("")
          localStorage.removeItem("logan_computer_number")
          this.navigate("/")
          throw new Error("세션이 만료되었습니다")
        } else if (response.status === 500) {
          alert("서버에 오류가 발생했습니다")
          throw new Error("서버에 오류 발생")
        } else {
          alert("오류가 발생했습니다")
          throw new Error(`오류가 발생했습니다 오류번호: ${response.status}`)
        }
      }

      const data = await response.json()

      return data

    } catch (err) {
      throw err
    }

  }

}