<template>
  <v-timeline :dense="$vuetify.breakpoint.sm || $vuetify.breakpoint.xs">
    <v-timeline-item
      v-for="(year, i) in years"
      :key="i"
      :color="year.color"
      large
    >
      <template v-slot:opposite>
        <div class="py-4">
          <div>
            Lorem ipsum dolor sit amet, no nam oblique veritus. Commune scaevola
            imperdiet nec ut, sed euismod convenire principes at. Est et nobis
            iisque percipit, an vim zril disputando voluptatibus, vix an
            salutandi sententiae.
          </div>
        </div>
      </template>
      <v-card class="mx-auto elevation-20" style="max-width: 400px">
        <v-row justify="space-between">
          <v-col sm="8" cols="12">
            <v-card-title>
              <div>
                <div :class="`${year.color}--text`" class="headline">
                  Lorem ipsum
                </div>
                <div class="subtitle-1">Ellie Goulding</div>
                <div class="subtitle-2">{{ year.year }}</div>
              </div>
            </v-card-title>
          </v-col>
          <v-col sm="4" cols="0">
            <v-img
              class="shrink"
              contain
              height="125px"
              src="https://picsum.photos/500/300?image=40"
              style="flex-basis: 125px"
            ></v-img>
          </v-col>
        </v-row>
        <v-divider dark></v-divider>
        <v-card-actions class="pa-4">
          Rate this
          <v-spacer></v-spacer>
          <span class="grey--text text--lighten-2 caption mr-2">
            ({{ 4 }})
          </span>
          <v-rating
            :value="4"
            background-color="white"
            color="yellow accent-4"
            dense
            half-increments
            hover
            size="18"
          ></v-rating>
        </v-card-actions>
      </v-card>
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