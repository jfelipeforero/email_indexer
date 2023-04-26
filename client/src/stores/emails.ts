import { defineStore } from 'pinia'
import axios from 'axios' 

export const useEmailsStore = defineStore('emails',{
   state: () => ({ 
    emails: [],
  }),

  actions: {
    async search(query:any){
      try {
        this.emails = await axios.get(`http://localhost:3000/search?q=${query.value}`) 
          console.log(this.emails)
      } catch (error) {
        console.error(error);
      } 
    }
  },
})
