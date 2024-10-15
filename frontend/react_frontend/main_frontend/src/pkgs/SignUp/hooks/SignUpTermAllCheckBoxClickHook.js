import { useCallback, useRef } from "react"

export const useSignUpTermAllCheckBoxClickHook = () => {
  const allCheckBoxDiv = useRef(null)

  const changeAllCheckBox = useCallback((event) => {
    if (event.target.className === "signUpTermAllCheckBoxInputValue") {
      const allCheckValue = event.target.checked
      allCheckBoxDiv.current = document.querySelectorAll(".signUpTermCheckBoxTermInputValue")
      Array.from(allCheckBoxDiv.current).forEach((checkBoxDiv) => {
        checkBoxDiv.checked = allCheckValue
      })
    }

  }, [])

  return { changeAllCheckBox }
}