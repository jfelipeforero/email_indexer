<template>
  <div class="table-container">
    <table>
      <thead>
        <tr>
          <th>Subject</th>
          <th>From</th>
          <th>To</th>
        </tr>
      </thead>
      <tbody>
        <tr>
          <td>{{emailsStore.emails[0]}}</td>
          <td>{{emailsStore.emails[0]}}</td>
          <td>{{emailsStore.emails[0]}}</td>
        </tr> 
      </tbody>
    </table>
  </div>
</template>

<script setup lang="ts">
import { useEmailsStore } from '../stores/emails'

const emailsStore = useEmailsStore();

</script>

<style>
.table-container {
  width: 50%; /* Set the width of the container to 50% */
  float: left; /* Float the container to the left */
  margin-left: 10px;
}

table { 
  border-collapse: collapse; /* Remove gaps between borders and cells */
  width: 100%; /* Set the width of the table to 100% */
}

th, td {
  border: 4px solid white; /* Set border of table cells */
  padding: 8px; /* Add padding to table cells */
  color:white;
  border-radius:10px;
}
</style>

