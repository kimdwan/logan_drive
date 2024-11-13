import "./assets/css/MainUser.css"

import { MainUserFooter, MainUserProfile } from "./components"
import { MainContext } from "../../../../App"

import { useContext } from "react"
export const MainUser = () => {
  const { computerNumber, setComputerNumber } = useContext(MainContext)

  return (
    <div className = "mainUserContainer">
      
      {/* 메인 프로필이 존재하는 컴퍼넌트 */}
      <MainUserProfile computerNumber = {computerNumber} setComputerNumber = {setComputerNumber} />

      {/* 메인 푸터 컴퍼넌트 */}
      <MainUserFooter />

    </div>
  )
}