import { useNavigate } from "react-router-dom"
import { useCallback, useRef } from "react"

export const useSignUpCheckOneTwoValueHook = () => {
    // 동의 버튼을 클릭했을때 확인해주는 훅
    const allCheckInputValue = useRef(null)
    const navigate = useNavigate()
    
    const clickGoSignUpFormBtn = useCallback((event) => {
      if (event.target.className === "signUpTermSubmitBtnValue") {
        try {
          allCheckInputValue.current = document.querySelectorAll(".signUpTermCheckBoxTermInputValue")
          let term_agree_3 = undefined
    
          Array.from(allCheckInputValue.current).forEach((checkBox, idx) => {
            const checkValue = checkBox.checked
            if (idx < 2) {
              if (!checkValue) {
                alert("필수 사항은 반드시 체크해 주셔야 합니다.")
                throw new Error("필수 사항 1,2를 체크하지 않았음")
              }
            } else {
              term_agree_3 = checkValue
            }
          })
  
          const url = `/signup/term=${term_agree_3}/form/`
          navigate(url)
  
        } catch (err) {
          throw err
        }
      }
  
    }, [ navigate ])

    return { clickGoSignUpFormBtn }
}