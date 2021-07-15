<template>
  <v-container class="fill-height">
    <v-row justify="center"> </v-row>
  </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { ListProject } from "../utils/api";
@Component({
  components: {},
})
export default class Project extends Vue {
  created() {
    ListProject({ Page: 0, PageSize: 1, Uid: "6047692a3ffb994a17160180" });
  }
}
</script>
<style>
.bordered {
  border: 2px solid #3298dc;
}
</style>
