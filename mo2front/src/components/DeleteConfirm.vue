<template>
  <v-row justify="center">
    <v-dialog v-model="show" persistent max-width="290" @click:outside="close">
      <v-card>
        <v-card-title class="headline">
          {{ title }}
        </v-card-title>
        <v-card-text>
          {{ content }}
        </v-card-text>
        <v-card-actions>
          <v-spacer></v-spacer>
          <v-btn color="red darken-1" text @click="confirm"> 确认删除 </v-btn>
          <v-btn color="green darken-1" text @click="close"> 取消 </v-btn>
        </v-card-actions>
      </v-card>
    </v-dialog>
  </v-row>
</template>
<script lang="ts">
import Vue from "vue";
import Component from "vue-class-component";
import { Prop } from "vue-property-decorator";
@Component({
  components: {},
})
export default class DeleteConfirm extends Vue {
  @Prop()
  title: string;
  @Prop()
  content: string;
  @Prop()
  show: boolean;
  confirm() {
    this.$emit("confirm");
    this.close();
  }
  close() {
    this.$emit("update:show", false);
  }
}
</script>