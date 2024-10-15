import "./assets/css/Main.css"
import { useContext } from "react"
import { MainContext } from "../../App"
import { MainUser, MainLogin, MainNavbar } from "./components"

export const Main = () => {
  const { computerNumber, setComputerNumber } = useContext(MainContext)

  return (
    <div className = "mainContainer">

      {/* 메인 네브바*/}
      <MainNavbar />

      {/* 메인 컴퍼넌트가 보임 */}
      {
        computerNumber ? <MainUser /> : <MainLogin setComputerNumber = {setComputerNumber} />
      }
    </div>
  )
}