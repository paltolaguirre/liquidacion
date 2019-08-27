DO $$
BEGIN		
	INSERT INTO liquidaciontipo(id,nombre, codigo, activo,created_at)
	VALUES 
		(-1,'Mensual','MENSUAL',1,current_timestamp),
		(-2,'Primer Quincena','PRIMER_QUINCENA',1,current_timestamp),
		(-3,'Segunda Quincena','SEGUNDA_QUINCENA',1,current_timestamp),
		(-4,'Vacaciones','VACACIONES',1,current_timestamp),
		(-5,'SAC','SAC',1,current_timestamp),
		(-6,'Liquidación Final','LIQUIDACION_FINAL',1,current_timestamp);

	INSERT INTO liquidacioncondicionpago(id,nombre,codigo,activo,created_at)
	VALUES
		(-1,'Cuenta Corriente','CUENTA_CORRIENTE',1,current_timestamp),
		(-2,'Contado','CONTADO',1,current_timestamp);
	
END;
$$;

