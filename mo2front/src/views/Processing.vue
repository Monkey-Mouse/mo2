<template>
  <v-container class="fill-height">
    <v-row justify="center">
      <v-col class="text-center">
        <v-progress-circular
          size="128"
          indeterminate
          color="primary"
        ></v-progress-circular>
      </v-col>
    </v-row>
    <v-row justify="center">
      <v-col class="text-center text-h3"> Processing </v-col>
    </v-row>
  </v-container>
</template>
<script lang="ts">
import axios, { AxiosError } from "axios";
import Vue from "vue";
import Component from "vue-class-component";
@Component({
  components: {},
})
export default class Processing extends Vue {
  created() {
    axios
      .get(this.$route.fullPath, { maxRedirects: 0 })
      .then((re) => {
        this.$router.push(
          (re.request.responseURL as string).replace('http://','https://').substring(
            window.location.host.length + window.location.protocol.length + 2
          )
        );
      })
      .catch((err: AxiosError) => {
        this.$router.push(
          (err.request.responseURL as string).replace('http://','https://').substring(
            window.location.host.length + window.location.protocol.length + 2
          )
        );
      });
  }
}
</script>