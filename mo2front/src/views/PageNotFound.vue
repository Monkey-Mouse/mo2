<template>
  <v-container class="fill-height">
    <v-row justify="center">
      <v-col class="text-center">
        <v-icon size="200">mdi-circle-off-outline</v-icon>
      </v-col>
      <v-col class="text-center text-h1"> 404 Page Not Found! </v-col>
    </v-row>
    <v-row justify="center">
      <a @click="$router.push('/')">Back to home</a>
    </v-row>
  </v-container>
</template>
<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
@Component({
  components: {},
})
export default class NotFound extends Vue {}
</script>