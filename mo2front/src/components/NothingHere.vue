<template>
  <v-container>
    <v-row justify="center">
      <v-col class="text-center"
        ><v-icon size="128">mdi-clipboard-text-off-outline</v-icon></v-col
      >
    </v-row>
    <v-row justify="center">
      <v-col class="text-center"
        ><h1 class="text-h1">
          Nothing here yet :<span v-if="showI">(</span><span v-else>)</span>
        </h1>
      </v-col>
    </v-row>
    <v-row justify="center">
      <v-col class="text-center"
        ><v-icon size="128">mdi-arrow-down-bold</v-icon></v-col
      >
    </v-row>
    <v-row justify="center">
      <v-col class="text-center"
        ><v-btn outlined @click="$emit('click')" color="primary">
          {{ btnText }}
        </v-btn>
      </v-col>
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
export default class NothingHere extends Vue {
  @Prop()
  btnText: string;
  showI = false;
  created() {
    setInterval(() => (this.showI = !this.showI), 500);
  }
}
</script>