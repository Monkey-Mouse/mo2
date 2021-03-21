<template>
  <v-container class="fill-height"> </v-container>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import BlogSkeleton from "../components/BlogTimeLineSkeleton.vue";
@Component({
  components: {},
})
export default class Notifications extends Vue {}
</script>
<style>
.bordered {
  border: 2px solid #3298dc;
}
</style>
