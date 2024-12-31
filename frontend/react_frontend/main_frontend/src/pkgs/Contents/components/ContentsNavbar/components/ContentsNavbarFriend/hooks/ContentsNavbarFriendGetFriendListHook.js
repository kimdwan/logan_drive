import { useEffect, useState } from "react"
import { ContentNavbarFriendFunc } from "../functions"

export const useContentNavbarFriendListHook = (computerNumber, setComputerNumber, navigate) => {
  const [ friendList, setFriendList ] = useState([])

  useEffect(() => {
    const getFriendListClass = new ContentNavbarFriendFunc(computerNumber, setComputerNumber, navigate)
    const getFriendData = async ( setFriendList ) => {

      try {
        const friendData = await getFriendListClass.getUserDetail()

        if (friendData && Array.from(friendData).length > 0) {
          
          setFriendList(friendData)

        }
      } catch (err) {
        throw err
      }

    }

    if (computerNumber) {
      getFriendData(setFriendList)
    }
  
  }, [ computerNumber, setComputerNumber, navigate ])

  return { friendList }
}