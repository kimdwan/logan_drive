

export class ContentsNavbarUploadImgFunction {
  constructor (computerNumber, setComputerNumber, navigate) {
    this.computerNumber = computerNumber
    this.setComputerNumber = setComputerNumber
    this.navigate = navigate
  }

  async UploadUserProfile(url, formData) {
    try {
      const response = await fetch(url, {
        method : "POST",
        headers : {
          "User-Computer-Number" : this.computerNumber
        },
        body : formData,
        credentials : "include",
      })

      if (!response.ok) {
        if (response.status === 401) {
          alert("세션이 만료되었습니다")
          localStorage.removeItem("logan_computer_number")
          this.setComputerNumber("")
          this.navigate("/")
          throw new Error("세션이 만료됨")
        } else if (response.status === 400) {
          alert("클라이언트에서 보낸 폼이 문제가 있습니다")
          throw new Error("클라이언트에서 보낸 폼에 오류가 있음")
        } else if (response.status === 500) {
          alert("서버에 오류가 발생했습니다")
          throw new Error("서버에 오류가 발생했습니다")
        } else {
          alert("오류가 발생했습니다")
          throw new Error(`오류 발생 오류 번호: ${response.status}`)
        }
      }

      const data = await response.json()
      
      if (data && data["message"]) {
        return data["message"]
      }


    } catch (err) {
      throw err
    }
  }
}