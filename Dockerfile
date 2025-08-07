FROM public.ecr.aws/lambda/provided:al2

# o provided:al2 espera o bootstrap em /var/runtime
COPY bootstrap ${LAMBDA_RUNTIME_DIR}/

# garante que é executável (opcional se já veio com +x)
RUN chmod +x ${LAMBDA_RUNTIME_DIR}/bootstrap

CMD ["bootstrap"]
