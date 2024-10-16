import * as yup from "yup"
import { useForm } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"

export const useSignUpFormHook = () => {
  const schema = yup.object({
    email : yup.string().email("이메일 형식을 지켜주세요").required("이메일은 필수로 입력해야 하는 사항입니다"),
    nickname : yup.string().min(3, "닉네임은 최소3글자 입니다").max(12, "닉네임은 최대12글자 입니다").required("닉네임은 필수로 입력해야 합니다"),
    password : yup.string().min(4, "비밀번호는 최소4글자 입니다").max(16, "비밀번호는 최대16글자 입니다").required("비밀번호는 필수로 입력해야 합니다."),
    confirm_password : yup.string().oneOf([yup.ref("password")], "비밀번호가 서로 다릅니다")
  })

  const { register, handleSubmit, formState:{ errors }, setError} = useForm({
    resolver : yupResolver(schema)
  })

  const onSubmit = (data) => {
    console.log(data)
  }



  return { register, handleSubmit, errors, onSubmit }
}