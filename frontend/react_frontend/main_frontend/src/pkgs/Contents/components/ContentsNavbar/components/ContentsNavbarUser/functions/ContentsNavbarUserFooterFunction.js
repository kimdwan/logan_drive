

export class ContentsNavbarUserFooterFunc{
  constructor(computerNumber, setComputerNumber, navigate) {
    this.computerNumber = computerNumber
    this.setComputerNumber = setComputerNumber
    this.navigate = navigate
  }

  async ClickLogoutBtn(url) {

    try {
      const response = await fetch(url, {
        method : "GET",
        headers : {
          "Content-Type" : "application/json; charset=utf-8",
          "User-Computer-Number" : this.computerNumber,
          "X-Requested-With" : "XMLHttpRequest",
        },
        credentials : "include",
      })

      if (!response.ok) {
        if (response.status === 401) {
          alert("세션이 만료되었습니다")
          localStorage.removeItem("logan_computer_number")
          this.setComputerNumber("")
          this.navigate("/")
          throw new Error("세션이 만료되었습니다")
        } else if (response.status === 500) {
          alert("서버에 오류가 발생했습니다")
          throw new Error("서버에 오류가 발생했습니다")
        } else {
          alert("오류가 발생했습니다")
          throw new Error(`오류 발생: ${response.status}`)
        }
      }

      const data = await response.json()

      if (data) {
        alert(data["message"])
        localStorage.removeItem("logan_computer_number")
        this.setComputerNumber("")
        this.navigate("/")
        return
      }


    } catch (err) {
      throw err
    }

  }

} 