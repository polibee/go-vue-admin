<script setup lang="ts">
import ResourceFormView from '@/components/resource/ResourceFormView.vue';
import { resourceApi } from '../api';
import { resourceDefinition } from '../resource';
</script>

<!-- Generated resource form. Register the route and matching backend contract after review. -->
<template>
  <ResourceFormView :resource="resourceDefinition" :api="resourceApi" />
</template>
