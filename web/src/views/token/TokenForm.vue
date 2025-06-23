<script setup lang="ts">
import { Coins, Info } from 'lucide-vue-next'
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { BaseInfoForm, FormActions, FormCard, PageHeader } from '@/components'
import { useSecretItemForm } from '@/composables/useSecretItemForm'
import { SECRET_ITEM_TYPE } from '@/constants'
import TokenDataForm from '@/views/token/TokenDataForm.vue'

const router = useRouter()
const route = useRoute()

// 是否为编辑模式
const isEdit = computed(() => !!route.params.id)

const {
  formData,
  loading,
  handleSubmit,
} = useSecretItemForm(SECRET_ITEM_TYPE.Token)
</script>

<template>
  <div class="space-y-6">
    <PageHeader
      :title="isEdit ? '编辑令牌' : '新建令牌'" :description="isEdit ? '修改现有的令牌信息' : '创建一个新的令牌'"
      @back="router.back()"
    />

    <form @submit="handleSubmit">
      <div class="grid gap-6 lg:grid-cols-2">
        <!-- 基础信息 -->
        <FormCard title="基础信息" description="填写令牌的基本信息" :icon="Info">
          <BaseInfoForm v-model="formData" :type="SECRET_ITEM_TYPE.Token" />
        </FormCard>

        <!-- 令牌信息 -->
        <FormCard title="令牌信息" description="配置令牌的具体信息" :icon="Coins">
          <TokenDataForm />
        </FormCard>
      </div>

      <FormActions :loading="loading" :is-edit="isEdit" @cancel="router.back()" />
    </form>
  </div>
</template>
