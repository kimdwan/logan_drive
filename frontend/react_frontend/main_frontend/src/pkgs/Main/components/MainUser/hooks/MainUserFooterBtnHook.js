import { MainUserFooterClass } from "../functions"

import { useCallback } from "react"

export const useMainUserFooterBtnHook = (computerNumber, setComputerNumber, go_backend_url ,navigate) => {
  const clickFooterBtn = useCallback(async (event) => {
    const footerClassFunc = new MainUserFooterClass(computerNumber, setComputerNumber, navigate)

    if (event.target.className === "mainUserFooterLogoutBtn") {

      const url = `${go_backend_url}/auth/user/logout`
      const message = await footerClassFunc.Logout(url)
      if (message) {
        alert(message)
        localStorage.removeItem("logan_computer_number")
        setComputerNumber("")
        navigate("/")
      }
    } else if ( event.target.className === "mainUserFooterGoMainBtn" ) {

      if (computerNumber) {
        footerClassFunc.GoMainContents()
      }
    }

  }, [ computerNumber, setComputerNumber, navigate, go_backend_url ])

  return { clickFooterBtn }
}