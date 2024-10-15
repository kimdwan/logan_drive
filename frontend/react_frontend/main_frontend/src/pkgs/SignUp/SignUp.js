import "./assets/css/SignUp.css"

import { SignUpForm, SignUpNavbar, SignUpTerm } from "./components"
import { useSignUpUrlTypeHook } from "./hooks"

export const SignUp = () => {
  const { urlPathType } = useSignUpUrlTypeHook()

  return (
    <div className = "signUpContainer">
      {/* 회원가입에서 네브바에 해당하는 컴퍼넌트 */}
      <SignUpNavbar />

      {/* 회원가입에 종착지를 설정하는 컴퍼넌트 */}
      {
        urlPathType === "term" ? <SignUpTerm /> : urlPathType === "form" ? <SignUpForm /> : <></>
      }
    </div>
  )
}