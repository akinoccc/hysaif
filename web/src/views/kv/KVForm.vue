<script setup lang="ts">
import { Braces, Info } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { BaseInfoForm, FormActions, FormCard, PageHeader } from '@/components'
import { useSecretItemForm } from '@/composables/useSecretItemForm'
import { SECRET_ITEM_TYPE } from '@/constants'
import CustomDataForm from '@/views/kv/KVDataForm.vue'

const router = useRouter()
const route = useRoute()

// 是否为编辑模式
const isEdit = computed(() => !!route.params.id)

const {
  form,
  formData,
  loading,
  handleSubmit,
} = useSecretItemForm(SECRET_ITEM_TYPE.KV)
</script>

<template>
  <div class="space-y-6">
    <PageHeader
      :title="isEdit ? '编辑 KV 键值对' : '新建 KV 键值对'" :description="isEdit ? '修改现有的 KV 键值对信息' : '创建一个新的 KV 键值对'"
      @back="router.back()"
    />

    <form @submit="handleSubmit">
      <div class="grid gap-6 lg:grid-cols-2">
        <!-- 基础信息 -->
        <FormCard title="基础信息" description="填写 KV 键值对的基本信息" :icon="Info">
          <BaseInfoForm v-model="formData" :type="SECRET_ITEM_TYPE.KV" />
        </FormCard>

        <!-- KV 键值对信息 -->
        <FormCard title="KV 键值对信息" description="配置 KV 键值对的具体信息" :icon="Braces">
          <CustomDataForm :form="form" />
        </FormCard>
      </div>

      <FormActions :loading="loading" :is-edit="isEdit" @cancel="router.back()" />
    </form>
  </div>
</template>
