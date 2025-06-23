<script setup lang="ts">
import { SecretList } from '@/components'
import { SECRET_ITEM_TYPE } from '@/constants'
</script>

<template>
  <SecretList :secret-type="SECRET_ITEM_TYPE.Token" />
</template>
