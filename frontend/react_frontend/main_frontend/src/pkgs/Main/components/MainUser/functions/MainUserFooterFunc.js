

export class MainUserFooterClass {

  constructor(computerNumber, setComputerNumber, navigate) {
    this.computerNumber = computerNumber
    this.setComputerNumber = setComputerNumber
    this.navigate = navigate
  }

  // 유저가 로그아웃 하게 해주는 함수
  async Logout(url) {

    try {
      
      const response = await fetch(url, {
        method : "GET",
        headers : {
          "Content-Type" : "application/json",
          "X-Requested-With" : "XMLHttpRequest",
          "User-Computer-Number" : this.computerNumber
        },
        credentials : "include",
      })

      if (!response.ok) {
        if (response.status === 401) {
          alert("세션이 만료되었습니다")
          localStorage.removeItem("logan_computer_number")
          this.setComputerNumber("")
          this.navigate("/")
          throw new Error("세션 만료")
        } else if (response.status === 500) {
          alert("서버에 오류가 발생했습니다")
          throw new Error("서버에 오류가 발생했습니다")
        } else {
          alert("오류가 발생했습니다")
          throw new Error(`오류가 발생했습니다 오류 번호: ${response.status}`)
        }
      }

      const data = await response.json()

      if (data) {
        return data["message"]
      }

    } catch (err) {
      throw err
    }

  }
  
  // 메인 컨텐츠를 이용 할 수 있게 해주는 함수
  GoMainContents() {
    if (this.computerNumber) {
      this.navigate("/main/channellist")
    }
  }

}