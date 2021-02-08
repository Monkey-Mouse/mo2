<template>
  <v-timeline>
    <v-timeline-item
      v-for="(year, i) in years"
      :key="i"
      :color="year.color"
      large
    >
      <template v-slot:opposite>
        <span
          :class="`headline font-weight-bold ${year.color}--text`"
          v-text="year.year"
        ></span>
        <span>
          <v-img
            lazy-src="https://picsum.photos/id/11/10/6"
            max-height="176"
            max-width="233"
            src="https://picsum.photos/500/300?image=40"
          ></v-img>
        </span>
      </template>
      <div class="py-4">
        <h2 :class="`headline font-weight-light mb-4 ${year.color}--text`">
          Lorem ipsum
        </h2>
        <div>
          Lorem ipsum dolor sit amet, no nam oblique veritus. Commune scaevola
          imperdiet nec ut, sed euismod convenire principes at. Est et nobis
          iisque percipit, an vim zril disputando voluptatibus, vix an salutandi
          sententiae.
        </div>
      </div>
    </v-timeline-item>
  </v-timeline>
</template>

<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
@Component
export default class Home extends Vue {
  years = [
    {
      color: "cyan",
      year: "1960",
    },
    {
      color: "green",
      year: "1970",
    },
    {
      color: "pink",
      year: "1980",
    },
    {
      color: "amber",
      year: "1990",
    },
    {
      color: "orange",
      year: "2000",
    },
  ];
}
</script>