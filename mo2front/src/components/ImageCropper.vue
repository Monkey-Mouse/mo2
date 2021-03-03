<template>
  <v-dialog :value="show" @click:outside="close" max-width="600px">
    <v-card>
      <v-container>
        <v-row>
          <v-col cols="12">
            <v-card-title> {{ title }} </v-card-title>
          </v-col>
        </v-row>
        <v-card-text>
          <div>
            <VueCropper
              style="height: 300px"
              ref="cropper"
              :img="img"
              :outputSize="1"
              :outputType="'webp'"
              :autoCrop="true"
              :fixed="true"
              :centerBox="true"
            ></VueCropper>
          </div>
        </v-card-text>
        <v-card-actions>
          <v-btn color="green" outlined text @click="confirm">确认</v-btn>
          <v-btn @click="close" color="red">取消</v-btn>
        </v-card-actions>
      </v-container>
    </v-card>
  </v-dialog>
</template>

<script lang="ts">
import { Prop } from "vue-property-decorator";
import Vue from "vue";
import Component from "vue-class-component";
import { VueCropper } from "vue-cropper";
@Component({
  components: { VueCropper },
})
export default class About extends Vue {
  @Prop({ default: "剪裁你的图片" })
  title: string;
  @Prop()
  show: boolean;
  @Prop()
  img: string | Blob;
  confirm() {
    // 获取截图的blob数据
    (this.$refs.cropper as any).getCropBlob((data) => {
      // do something
      this.$emit("confirm", data);
    });
  }
  close() {
    this.$emit("update:show", false);
  }
}
</script>

<style lang="scss" scoped>
.v-icon.v-icon.v-icon--link {
  padding: 0 10px;
}
</style>
