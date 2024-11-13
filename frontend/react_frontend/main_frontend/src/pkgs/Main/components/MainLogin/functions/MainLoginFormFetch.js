

export const MainLoginFormFetch = async ( url, datas, setError, setComputerNumber ) => {
  try {
    const response = await fetch(url, {
      method : "POST",
      headers : {
        "Content-Type" : "application/json",
        "X-Requested-With" : "XMLHttpRequest",
      },
      body :JSON.stringify(datas),
      credentials : "include",
    })

    if (!response.ok) {
      if (response.status === 406) {
        setError("email", {
          type : "manual",
          message : "이메일을 찾을수가 없습니다."
        })
        throw new Error("이메일을 찾을수가 없습니다")
      } else if (response.status === 510) {
        setError("password", {
          type : "manual",
          message : "비밀번호가 틀렸습니다."
        })
        throw new Error("비밀번호가 틀렸습니다")
      } else if (response.status === 400) {
        alert("클라이언트에서 보낸 폼에 문제가 있습니다")
        throw new Error("클라이언트에서 보낸 폼에 문제가 있습니다")
      } else if (response.status === 500) {
        alert("서버에 오류가 있습니다.")
        throw new Error("서버에 오류가 있습니다.")
      } else {
        alert("오류가 발생했습니다.") 
        throw new Error(`오류가 발생했습니다 오류번호: ${response.status}`)
      }
    }

    const data = await response.json()
    if (data && data["computer_number"]) {
      const computer_number = data["computer_number"]
      setComputerNumber(computer_number)
      localStorage.setItem("logan_computer_number",computer_number)
      return data["message"]
    }

  } catch (err) {
    throw err
  }
}