import * as yup from "yup"
import { useForm } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import { MainLoginFormFetch } from "../functions"

export const useMainLoginFormHook = (setComputerNumber) => {
  const schema = yup.object({
    email : yup.string().email("이메일 형식을 지켜주시길 바랍니다").required("이메일은 필수로 작성해야 합니다."),
    password : yup.string().min(6, "비밀번호는 최소 6글자 입니다.").max(16, "비밀번호는 최대 16글자 입니다").required("비밀번호는 필수로 입력해야 합니다.")
  })

  const { register, handleSubmit, formState : { errors }, setError} = useForm({
    resolver : yupResolver(schema)
  })

  const onSubmit = async (data) => {

    try {
      // 비밀번호 최소검증
      const passwordValue = data.password
      const passwordRegex = /^(?=.*[A-Za-z])(?=.*\d)(?=.*[!@#$%^&*])[A-Za-z\d!@#$%^&*]+$/
      if (!passwordRegex.test(passwordValue)) {
        setError("password", {
          type : "manual",
          message : "비밀번호는 문자, 숫자, 특수문자 한개 이상씩 들어가 있는 6글자 문자 이어야 합니다."
        })
        throw new Error("비밀번호 형식을 지키지 않았음")
      }

      if (data) {
        const go_backend_url = process.env.REACT_APP_GO_BACKEND_URL
        const url = `${go_backend_url}/user/login`
        const message = await MainLoginFormFetch(url, data, setError, setComputerNumber)
        if (message) {
          alert(message)
        }
      }
 
    } catch(err) {
      throw err
    }
  }

  return { register, handleSubmit, errors, onSubmit }
}