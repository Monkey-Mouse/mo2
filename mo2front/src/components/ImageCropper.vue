<template>
  <v-dialog :value="show" @click:outside="close" max-width="600px">
    <v-card :loading="loading" :disabled="loading">
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
              @imgLoad="imgLoad"
            ></VueCropper>
          </div>
          <v-row v-if="confirmerr !== ''">
            <v-alert dense outlined type="error" class="col-12">{{
              confirmerr
            }}</v-alert></v-row
          >
        </v-card-text>
        <v-card-actions>
          <v-btn color="success" outlined text @click="confirm">确认</v-btn>
          <v-btn @click="close" color="error">取消</v-btn>
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
  @Prop({ default: false })
  loading: boolean;
  @Prop({ default: "" })
  confirmerr: string;
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
  imgLoad(success) {
    this.$emit("imgLoad", success);
  }
}
</script>

<style lang="scss" scoped>
.v-icon.v-icon.v-icon--link {
  padding: 0 10px;
}
</style>
