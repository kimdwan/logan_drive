import { useEffect, useState } from "react"

export const useSignUpGetTermAgree3Hook = () => {
  const [ termAgree3, setTermAgree3 ] = useState(undefined)

  useEffect(() => {
    const current_url = window.location.href
    const current_url_list = current_url.split("/")
    const term_agree_name = current_url_list[current_url_list.length - 3]
    const term_agree_value = term_agree_name.split("=")
    const term_agree_3 = term_agree_value[1]
    if (term_agree_3 === "true") {
      setTermAgree3(true)
    } else if (term_agree_3 === "false") {
      setTermAgree3(false)
    }
  },[ ])

  return { termAgree3 }
}