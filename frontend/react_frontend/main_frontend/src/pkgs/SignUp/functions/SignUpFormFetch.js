

export const SignUpFormFetch = async ( url, datas, setError ) => {
  try {
    const response = await fetch(url, {
      method : "POST",
      headers : {
        "Content-Type" : "application/json",
        "X-Requested-With" : "XMLHttpRequest",
      },
      body : JSON.stringify(datas)
    })

    if (!response.ok) {
      if (response.status === 406) {
        setError("email", {
          type : "manual",
          message : "가입한 이메일이 존재합니다."
        })
        throw new Error("이메일 중복")
      } else if (response.status === 510) {
        setError("nickname", {
          type : "manual",
          message : "닉네임이 중복됩니다."
        })
        throw new Error("닉네임 중복")
      } else if (response.status === 400) {
        alert("클라이언트에서 보낸 값에 문제가 있습니다")
        throw new Error("클라이언트값 오류")
      } else if (response.status === 500) {
        alert("서버에 오류가 발생했습니다")
        throw new Error("서버에 오류 발생")
      } else {
        alert("오류가 발생했습니다")
        throw new Error(`오류 번호: ${response.status}`)
      }
    }

    const backend_data = await response.json()
    
    if (backend_data) {
      const message = backend_data["message"]
      return message
    }

  } catch (err) {
    throw err
  }
}