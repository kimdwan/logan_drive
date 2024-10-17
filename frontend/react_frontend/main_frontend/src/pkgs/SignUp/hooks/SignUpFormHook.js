import * as yup from "yup"
import { useForm } from "react-hook-form"
import { yupResolver } from "@hookform/resolvers/yup"
import { useNavigate } from "react-router-dom"

import { SignUpFormFetch } from "../functions"

export const useSignUpFormHook = (termAgree3) => {
  const schema = yup.object({
    email : yup.string().email("이메일 형식을 지켜주세요").required("이메일은 필수로 입력해야 하는 사항입니다"),
    nickname : yup.string().min(3, "닉네임은 최소3글자 입니다").max(12, "닉네임은 최대12글자 입니다").required("닉네임은 필수로 입력해야 합니다"),
    password : yup.string().min(4, "비밀번호는 최소4글자 입니다").max(16, "비밀번호는 최대16글자 입니다").required("비밀번호는 필수로 입력해야 합니다."),
    confirm_password : yup.string().oneOf([yup.ref("password")], "비밀번호가 서로 다릅니다")
  })

  const { register, handleSubmit, formState:{ errors }, setError} = useForm({
    resolver : yupResolver(schema)
  })

  const navigate = useNavigate()

  const onSubmit = async (data) => {

    try {
      // 비밀번호 형식이 제대로 진행되는지 확인하는 함수
      const checkPasswordText = /[A-Za-z]/
      const checkPasswordNumber = /[0-9]/
      const checkPasswordChar = /[^A-Za-z0-9]/
      const password = data.password

      if (!checkPasswordText.test(password) || !checkPasswordNumber.test(password) || !checkPasswordChar.test(password)) {
        setError("password", {
          type : "manual",
          message : "비밀번호에는 문자, 숫자, 특수문자를 포함하고 있어야 합니다",
        })
        throw new Error("비밀번호 형식에 문제가 있음")
      }

      // 유저의 체크표시 업데이트
      data["term_agree_3"] = termAgree3

      // 회원가입을 진행하는 함수
      const backend_url = process.env.REACT_APP_GO_BACKEND_URL
      const url = `${backend_url}/user/signup`

      const message = await SignUpFormFetch( url, data, setError )
      if (message) {
        alert(message)
        console.clear()
        navigate("/")
      }
    } catch (err) {
      throw err
    }
  }



  return { register, handleSubmit, errors, onSubmit }
}