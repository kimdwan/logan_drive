import "./assets/css/MainLogin.css"
import { MainLoginLogo, MainLoginForm } from "./components"

export const MainLogin = ({ setComputerNumber }) => {
  return (
    <div className = "mainLoginContainer">

      {/* 로고가 들어오는 컴퍼넌트 */}
      <MainLoginLogo />

      {/* 로그인이 이루어지는 컴퍼넌트 */}
      <MainLoginForm setComputerNumber = {setComputerNumber} />

    </div>
  )
}