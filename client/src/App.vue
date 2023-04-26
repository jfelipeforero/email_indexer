<template>
  <div class= "full-page" >
  <SearchBar/>   
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue';
import TableComponent from './components/TableComponent.vue';
import SearchBar from './components/SearchBar.vue' 

export default defineComponent({
  components: {
    SearchBar,
    TableComponent
  }
});
</script>

<style>
.full-page{
  width: 100%;
  height: 100%;

}
</style>
