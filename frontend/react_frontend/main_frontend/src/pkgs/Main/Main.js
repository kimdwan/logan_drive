import "./assets/css/Main.css"
import { useContext } from "react"
import { MainContext } from "../../App"
import { MainUser, MainLogin } from "./components"

export const Main = () => {
  const { computerNumber, setComputerNumber } = useContext(MainContext)

  return (
    <div className = "mainContainer">
      {
        computerNumber ? <MainUser /> : <MainLogin setComputerNumber = {setComputerNumber} />
      }
    </div>
  )
}