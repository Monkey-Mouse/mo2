<template>
  <v-timeline :dense="this.$vuetify.breakpoint.smAndDown">
    <v-timeline-item
      v-for="i in num"
      :key="i"
      :color="displayColors[i % displayColors.length]"
      large
    >
      <v-skeleton-loader type="card"></v-skeleton-loader>
      <template v-slot:opposite>
        <v-skeleton-loader type="paragraph"></v-skeleton-loader>
      </template>
    </v-timeline-item>
  </v-timeline>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
import { colors } from "vuetify/lib";
@Component
export default class BlogTimeLineSkeleton extends Vue {
  @Prop()
  num: number;
  displayColors: string[] = [];
  created() {
    this.displayColors = Object.getOwnPropertyNames(colors);
    // let count = Object.getOwnPropertyNames(colors).length;
    // while (this.displayColors.length < count) {
    //   let copiedColors: Colors = Copy(colors) as Colors;
    //   let displayColor = randomProperty(copiedColors) as string;
    //   delete copiedColors[displayColor];
    //   this.displayColors.push(displayColor);
    // }
  }
}
</script>
<style>
.unclickable {
  cursor: default;
}
</style>