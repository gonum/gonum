#!/usr/bin/env ruby
# Generate vocab.jsonld and vocab.html from vocab.ttl and vocab_template.
#
# Generating vocab.jsonld is equivalent to running the following:
#
#    jsonld --compact --context vocab_context.jsonld --input-format ttl vocab.ttl  -o vocab.jsonld
require 'linkeddata'
require 'haml'
require 'active_support'

File.open("vocab.jsonld", "w") do |f|
  r = RDF::Repository.load("vocab.ttl")
  JSON::LD::API.fromRDF(r, useNativeTypes: true) do |expanded|
    # Remove leading/trailing and multiple whitespace from rdf:comments
    expanded.each do |o|
      c = o[RDF::RDFS.comment.to_s].first['@value']
      o[RDF::RDFS.comment.to_s].first['@value'] = c.strip.gsub(/\s+/m, ' ')
    end
    JSON::LD::API.compact(expanded, File.open("vocab_context.jsonld")) do |compacted|
      # Create vocab.jsonld
      f.write(compacted.to_json(JSON::LD::JSON_STATE))

      # Create vocab.html using vocab_template.haml and compacted vocabulary
      template = File.read("vocab_template.haml")
      
      html = Haml::Engine.new(template, :format => :html5).render(self,
        ontology:   compacted['@graph'].detect {|o| o['@id'] == "http://json-ld.github.io/normalization/tests/vocab#"},
        classes:    compacted['@graph'].select {|o| o['@type'] == "rdfs:Class"}.sort_by {|o| o['rdfs:label']},
        properties: compacted['@graph'].select {|o| o['@type'] == "rdf:Property"}.sort_by {|o| o['rdfs:label']}
      )
      File.open("vocab.html", "w") {|fh| fh.write html}
    end
  end
end
