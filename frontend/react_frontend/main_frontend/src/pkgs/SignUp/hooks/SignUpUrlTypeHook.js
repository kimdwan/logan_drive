import { useEffect, useState } from "react"
import { useLocation } from "react-router-dom"

export const useSignUpUrlTypeHook = () => {
  const [ urlPathType, setUrlPathType ] = useState("term")
  const urlLocation = useLocation()
  
  useEffect(() => {
    const urlName = urlLocation.pathname
    const urlNameList = urlName.split("/")
    const pickUrlName = urlNameList[urlNameList.length - 2]
    setUrlPathType(pickUrlName)
  }, [ urlLocation ])

  return { urlPathType }
}